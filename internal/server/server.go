package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/jonwraymond/toolprotocol/a2a"
)

// Config configures the HTTP server.
type Config struct {
	Host              string
	Port              int
	BasePath          string
	ReadHeaderTimeout time.Duration
}

// Server hosts the A2A HTTP endpoints.
type Server struct {
	cfg    Config
	handle *a2a.Handler
	server *http.Server
}

// New creates a new server.
func New(cfg Config, handler *a2a.Handler) *Server {
	return &Server{
		cfg:    cfg,
		handle: handler,
	}
}

// Run starts the HTTP server and blocks until context cancellation.
func (s *Server) Run(ctx context.Context) error {
	host := s.cfg.Host
	if host == "" {
		host = "0.0.0.0"
	}
	port := s.cfg.Port
	if port == 0 {
		port = 8091
	}
	base := s.cfg.BasePath
	if base == "" {
		base = "/a2a"
	}

	mux := http.NewServeMux()
	mux.HandleFunc(base, s.handle.ServeRPC)
	mux.HandleFunc(base+"/agent-card", s.handle.ServeAgentCard)
	mux.HandleFunc(base+"/skills", s.handle.ServeSkills)
	mux.HandleFunc(base+"/tasks", s.handle.ServeTaskList)
	mux.HandleFunc(base+"/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, base+"/tasks/")
		parts := strings.Split(path, "/")
		taskID := parts[0]
		if taskID == "" {
			http.NotFound(w, r)
			return
		}
		if len(parts) == 1 {
			s.handle.ServeTask(w, r, taskID)
			return
		}
		if len(parts) == 2 && parts[1] == "events" {
			s.handle.ServeTaskEvents(w, r, taskID)
			return
		}
		http.NotFound(w, r)
	})

	addr := fmt.Sprintf("%s:%d", host, port)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: s.cfg.ReadHeaderTimeout,
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}

	s.server = httpServer

	errCh := make(chan error, 1)
	go func() {
		err := httpServer.Serve(ln)
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		_ = s.Close()
		return nil
	case err := <-errCh:
		if err == nil {
			return nil
		}
		return err
	}
}

// Close stops the HTTP server.
func (s *Server) Close() error {
	if s.server == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}
