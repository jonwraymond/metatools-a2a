package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jonwraymond/metatools-a2a/internal/agent"
	"github.com/jonwraymond/metatools-a2a/internal/config"
	"github.com/jonwraymond/metatools-a2a/internal/server"
	"github.com/jonwraymond/tooldiscovery/discovery"
	"github.com/jonwraymond/tooldiscovery/tooldoc"
	"github.com/jonwraymond/toolexec/run"
	"github.com/jonwraymond/toolfoundation/model"
	"github.com/jonwraymond/toolprotocol/a2a"
	"github.com/jonwraymond/toolprotocol/task"
	"gopkg.in/yaml.v3"
)

type bootstrapFile struct {
	Tools []toolRegistration `yaml:"tools"`
}

type toolRegistration struct {
	Tool    model.Tool        `yaml:"tool"`
	Backend model.ToolBackend `yaml:"backend"`
	Doc     *tooldoc.DocEntry `yaml:"doc,omitempty"`
}

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	disc, err := discovery.New(discovery.Options{})
	if err != nil {
		log.Fatalf("discovery error: %v", err)
	}

	if cfg.Bootstrap.ToolsFile != "" {
		if err := loadBootstrap(cfg.Bootstrap.ToolsFile, disc); err != nil {
			log.Fatalf("bootstrap error: %v", err)
		}
	}

	runner := run.NewRunner(run.WithIndex(disc.Index()))

	baseURL := fmt.Sprintf("http://%s:%d%s", hostForURL(cfg.Server.Host), cfg.Server.Port, cfg.Server.BasePath)
	agentSvc := &agent.Agent{
		Name:             cfg.Provider.Name,
		Description:      cfg.Provider.Description,
		Version:          cfg.Provider.Version,
		DocumentationURL: cfg.Provider.DocumentationURL,
		IconURL:          cfg.Provider.IconURL,
		BaseURL:          baseURL,
		Discovery:        disc,
		Runner:           runner,
		MaxSkills:        cfg.Bootstrap.MaxSkills,
	}

	handler := a2a.NewHandler(agentSvc, task.NewManager())
	srv := server.New(server.Config{
		Host:     cfg.Server.Host,
		Port:     cfg.Server.Port,
		BasePath: cfg.Server.BasePath,
	}, handler)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := srv.Run(ctx); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func loadBootstrap(path string, disc *discovery.Discovery) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var file bootstrapFile
	if err := yaml.Unmarshal(data, &file); err != nil {
		return err
	}
	for _, reg := range file.Tools {
		if err := disc.RegisterTool(reg.Tool, reg.Backend, reg.Doc); err != nil {
			return err
		}
	}
	return nil
}

func hostForURL(host string) string {
	if host == "" || host == "0.0.0.0" {
		return "localhost"
	}
	return host
}
