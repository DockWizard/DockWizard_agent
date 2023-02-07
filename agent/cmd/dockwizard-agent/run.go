package main

import (
	"github.com/dockwizard/dockwizard/agent/pkg/agent"
	"github.com/dockwizard/dockwizard/agent/pkg/backend"
	"github.com/dockwizard/dockwizard/agent/pkg/backend/api"
	"github.com/dockwizard/dockwizard/agent/pkg/backend/stdout"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var run = &cobra.Command{
	Use:   "run",
	Short: "Run the agent",
	Run:   runMn,
}

func runMn(_ *cobra.Command, _ []string) {
	var b backend.Backend
	switch cfg.Backend {
	case "stdout":
		b = &stdout.Stdout{}
	case "api":
		b = api.New(cfg.APIEndpoint, cfg)
	}

	agentInstance, err := agent.New(cfg, b)
	if err != nil {
		logrus.Errorf("failed to create agent: %v", err)
	}
	agentInstance.Run()
}
