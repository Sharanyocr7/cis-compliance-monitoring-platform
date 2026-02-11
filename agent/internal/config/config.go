package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Agent struct {
		Name          string `yaml:"name"`
		Version       string `yaml:"version"`
		Build         string `yaml:"build"`
		SchemaVersion string `yaml:"schema_version"`
	} `yaml:"agent"`

	Backend struct {
		URL    string `yaml:"url"`
		APIKey string `yaml:"api_key"`
	} `yaml:"backend"`

	Runtime struct {
		TimeoutSeconds int `yaml:"timeout_seconds"`
		RetryCount     int `yaml:"retry_count"`
		RetryBackoffMs int `yaml:"retry_backoff_ms"`
	} `yaml:"runtime"`
}

func Load(path string) (Config, error) {
	var cfg Config
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}

	// sensible defaults
	if cfg.Agent.Name == "" {
		cfg.Agent.Name = "cis-agent"
	}
	if cfg.Agent.SchemaVersion == "" {
		cfg.Agent.SchemaVersion = "1.0"
	}
	if cfg.Runtime.TimeoutSeconds <= 0 {
		cfg.Runtime.TimeoutSeconds = 25
	}
	if cfg.Runtime.RetryCount < 0 {
		cfg.Runtime.RetryCount = 0
	}
	if cfg.Runtime.RetryBackoffMs <= 0 {
		cfg.Runtime.RetryBackoffMs = 800
	}
	return cfg, nil
}

func (c Config) Timeout() time.Duration {
	return time.Duration(c.Runtime.TimeoutSeconds) * time.Second
}
