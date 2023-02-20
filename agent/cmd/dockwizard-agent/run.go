package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/docker/docker/client"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/agent"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/backend"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/backend/api"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/backend/stdout"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
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
		b = stdout.New()
	case "api":
		b = api.New(cfg.APIEndpoint, cfg, nil)
	}

	ctx, cancel := context.WithCancel(context.Background())

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		cancel()
		log.Fatalf("failed to create docker client: %v", err)
	}

	agentInstance := agent.New(cfg, b, cli)
	if err != nil {
		logrus.Errorf("failed to create agent: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("gracefully shutting down")
			cancel()
		}
	}()

	agentInstance.Run(ctx)
}
