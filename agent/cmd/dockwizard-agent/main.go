package main

import (
	"github.com/dockwizard/dockwizard_agent/agent/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfg = &config.Config{}

var root = &cobra.Command{
	Use:              "dockwizard-agent",
	PersistentPreRun: preRun,
}

func init() {
	root.PersistentFlags().String("config-file", "/etc/dockwizard/dockwizard.yaml", "Path to the configuration file")

	root.AddCommand(run)
}

func preRun(cmd *cobra.Command, args []string) {
	path, err := cmd.Flags().GetString("config-file")
	if err != nil {
		logrus.Fatal(err)
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
