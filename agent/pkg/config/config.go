package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	// Agent ID is the unique ID of the agent
	AgentID string `yaml:"agent_id"`

	// Backend is the endpoint to send data to
	// Can also be "stdout" to print to stdout
	Backend string `yaml:"backend"`

	// API endpoint is the endpoint to send data to
	// Only used if backend is "api"
	APIEndpoint string `yaml:"api_endpoint"`

	// PollInterval is the interval to poll for new data
	// The value is in seconds, minimum 2 seconds
	PollInterval int `yaml:"poll_interval"`
}

func Read(path string) (*Config, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	bts, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(bts, &cfg)
	if err != nil {
		return nil, err
	}

	// Validate the config
	if cfg.PollInterval < 2 {
		return nil, fmt.Errorf("poll interval must be at least 2 seconds")
	}
	if cfg.Backend == "" {
		return nil, fmt.Errorf("backend is required, use stdout to print to stdout")
	}

	return &cfg, nil
}
