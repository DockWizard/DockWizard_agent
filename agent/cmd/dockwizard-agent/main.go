package main

import (
	"os"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var cfg = &config.Config{}

var root = &cobra.Command{
	Use:              "dockwizard-agent",
	PersistentPreRun: preRun,
}

func init() {
	root.PersistentFlags().String("config-file", "/etc/dockwizard/dockwizard.yaml", "Path to the configuration file")
	root.PersistentFlags().String("api-key", "", "Setup the API key for the agent")
	root.PersistentFlags().String("api-endpoint", "", "Setup the API endpoint for the agent")
	root.PersistentFlags().Int("update-frequency", 2, "Setup the update frequency for the agent (in seconds)")
	root.PersistentFlags().Bool("setup", false, "Setup the agent")
	root.PersistentFlags().Bool("overwrite-config", false, "Overwrite the config file if it exists")

	root.AddCommand(run)
}

func preRun(cmd *cobra.Command, args []string) {
	setup, _ := cmd.Flags().GetBool("setup")
	path, err := cmd.Flags().GetString("config-file")
	if err != nil {
		logrus.Fatal(err)
	}

	// Setup the agent if setup is set
	if setup {
		// Require api-key, api-endpoint and update-frequency
		apiKey, err := cmd.Flags().GetString("api-key")
		if err != nil {
			logrus.Fatal(err)
		}
		apiEndpoint, err := cmd.Flags().GetString("api-endpoint")
		if err != nil {
			logrus.Fatal(err)
		}
		updateFrequency, err := cmd.Flags().GetInt("update-frequency")
		if err != nil {
			logrus.Fatal(err)
		}

		if apiKey == "" {
			logrus.Fatal("api-key is required")
		}
		if apiEndpoint == "" {
			logrus.Fatal("api-endpoint is required")
		}
		if updateFrequency < 2 {
			logrus.Fatal("update-frequency must be at least 2 seconds")
		}

		// Check if the config file exists
		_, err = os.Stat(path)
		if err == nil {
			// Config file exists, check if we should overwrite it
			overwrite, _ := cmd.Flags().GetBool("overwrite-config")
			if !overwrite {
				logrus.Fatalf("config file already exists, use --overwrite-config to overwrite it")
			}
		}

		// Create the config file
		cfg := &config.Config{
			APIKey:          apiKey,
			Backend:         "api",
			APIEndpoint:     apiEndpoint,
			UpdateFrequency: updateFrequency,
		}

		// Write the config file
		bts, err := yaml.Marshal(cfg)
		if err != nil {
			logrus.Fatalf("failed to marshal config: %v", err)
		}
		err = os.WriteFile(path, bts, 0644)
		if err != nil {
			logrus.Fatalf("failed to write config file: %v", err)
		}
	}

	cfg, err = config.Read(path)
	if err != nil {
		logrus.Fatalf("failed to read config file: %v", err)
	}
}

func main() {
	if err := root.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
