package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Container struct{}

type Config struct {
	// API Key is the key to use when sending data to the backend
	APIKey string `yaml:"api_key"`

	// Backend is the endpoint to send data to
	// Can also be "stdout" to print to stdout
	Backend string `yaml:"backend"`

	// API endpoint is the endpoint to send data to
	// Only used if backend is "api"
	APIEndpoint string `yaml:"api_endpoint"`

	// UpdateFrequency is the interval to poll for new data
	// The value is in seconds, minimum 2 seconds
	UpdateFrequency int `yaml:"update_frequency"`

	// Containers is a list of containers to apply health checks to
	Containers []Container `yaml:"containers"`
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
	if cfg.UpdateFrequency < 2 {
		return nil, fmt.Errorf("update frequency must be at least 2 seconds")
	}
	if cfg.Backend == "" {
		return nil, fmt.Errorf("backend is required, use stdout to print to stdout")
	}

	// Verify that API endpoint and API key are set if backend is "api"
	if cfg.Backend == "api" {
		if cfg.APIEndpoint == "" {
			return nil, fmt.Errorf("api endpoint is required when backend is api")
		}
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("api key is required when backend is api")
		}
	}

	return &cfg, nil
}
