package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds runtime configuration for the A2A server.
type Config struct {
	Server    ServerConfig    `yaml:"server"`
	Provider  ProviderConfig  `yaml:"provider"`
	Bootstrap BootstrapConfig `yaml:"bootstrap"`
}

// ServerConfig configures the HTTP server.
type ServerConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	BasePath string `yaml:"basePath"`
}

// ProviderConfig configures the agent card identity.
type ProviderConfig struct {
	Name             string `yaml:"name"`
	Description      string `yaml:"description"`
	Version          string `yaml:"version"`
	DocumentationURL string `yaml:"documentationUrl"`
	IconURL          string `yaml:"iconUrl"`
}

// BootstrapConfig controls tool bootstrap behavior.
type BootstrapConfig struct {
	ToolsFile string `yaml:"toolsFile"`
	MaxSkills int    `yaml:"maxSkills"`
}

// DefaultConfig returns the default configuration.
func DefaultConfig() Config {
	return Config{
		Server: ServerConfig{
			Host:     "0.0.0.0",
			Port:     8091,
			BasePath: "/a2a",
		},
		Provider: ProviderConfig{
			Name:        "metatools-a2a",
			Description: "A2A reference server for ApertureStack tools",
			Version:     "0.1.0",
		},
		Bootstrap: BootstrapConfig{
			MaxSkills: 500,
		},
	}
}

// Load reads configuration from a YAML file and environment variables.
// Environment variables override file values.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		path = os.Getenv("METATOOLS_A2A_CONFIG")
	}
	if path != "" {
		// #nosec G304 -- config path is explicitly user-supplied via CLI/env.
		data, err := os.ReadFile(path)
		if err != nil {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return Config{}, fmt.Errorf("parse config: %w", err)
		}
	}

	applyEnv(&cfg)
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if v := os.Getenv("METATOOLS_A2A_HOST"); v != "" {
		cfg.Server.Host = v
	}
	if v := os.Getenv("METATOOLS_A2A_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = port
		}
	}
	if v := os.Getenv("METATOOLS_A2A_BASE_PATH"); v != "" {
		cfg.Server.BasePath = v
	}
	if v := os.Getenv("METATOOLS_A2A_NAME"); v != "" {
		cfg.Provider.Name = v
	}
	if v := os.Getenv("METATOOLS_A2A_DESCRIPTION"); v != "" {
		cfg.Provider.Description = v
	}
	if v := os.Getenv("METATOOLS_A2A_VERSION"); v != "" {
		cfg.Provider.Version = v
	}
	if v := os.Getenv("METATOOLS_A2A_DOCS_URL"); v != "" {
		cfg.Provider.DocumentationURL = v
	}
	if v := os.Getenv("METATOOLS_A2A_ICON_URL"); v != "" {
		cfg.Provider.IconURL = v
	}
	if v := os.Getenv("METATOOLS_A2A_TOOLS_FILE"); v != "" {
		cfg.Bootstrap.ToolsFile = v
	}
	if v := os.Getenv("METATOOLS_A2A_MAX_SKILLS"); v != "" {
		if max, err := strconv.Atoi(v); err == nil {
			cfg.Bootstrap.MaxSkills = max
		}
	}
}
