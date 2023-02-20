package agent

import (
	"context"
	"io"
	"log"
	"math"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/backend"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/config"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/data"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/dockerstats"
	"github.com/sirupsen/logrus"
)

type Agent struct {
	Config *config.Config

	buffer  []*data.Metrics
	docker  client.APIClient
	backend backend.Backend
}

func New(c *config.Config, b backend.Backend, cli client.APIClient) *Agent {
	return &Agent{
		Config:  c,
		buffer:  []*data.Metrics{},
		docker:  cli,
		backend: b,
	}
}

func (a *Agent) getDockerContainerMetrics() ([]*data.ContainerMetrics, error) {
	var ret []*data.ContainerMetrics

	ctx := context.TODO()

	// Get the metrics
	allContainers, err := a.docker.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// Get the metrics for each container
	for _, container := range allContainers {
		stats, err := a.docker.ContainerStatsOneShot(ctx, container.ID)
		if err != nil {
			return nil, err
		}

		statsBytes, err := io.ReadAll(stats.Body)
		if err != nil {
			return nil, err
		}
		parsedStats, err := dockerstats.Unmarshal(statsBytes)
		if err != nil {
			return nil, err
		}

		rx, tx := parsedStats.NetworkStats()
		read, write := parsedStats.DiskStats()

		ret = append(ret, &data.ContainerMetrics{
			ID:                    container.ID,
			Name:                  strings.TrimPrefix(container.Names[0], "/"),
			Image:                 container.Image,
			CPUUsage:              math.Round(parsedStats.CpuUsagePercentage()*1000) / 1000,
			MemoryUsage:           parsedStats.UsedMemory(),
			MemoryUsagePercentage: math.Round(parsedStats.MemoryUsagePercentage()*1000) / 1000,
			State:                 container.State,
			NetworkIORead:         int(rx),
			NetworkIOWrite:        int(tx),
			BlockIORead:           int(read),
			BlockIOWrite:          int(write),
		})
	}

	return ret, nil
}

func (a *Agent) sleep() {
	time.Sleep(time.Duration(a.Config.UpdateFrequency) * time.Second)
}

func (a *Agent) Run(ctx context.Context) {
	for {
		// Sleep for the poll interval
		a.sleep()

		// If canceled then break
		if ctx.Err() == context.Canceled {
			break
		}

		containerMetrics, err := a.getDockerContainerMetrics()
		if err != nil {
			log.Printf("error getting container metrics: %v", err)
			continue
		}

		// Send the metrics to the backend
		metrics := &data.Metrics{
			Container: containerMetrics,
		}
		err = a.backend.SendData(metrics)
		if err != nil {
			log.Printf("could not send metrics to backend: %v", err)
		}
	}

	err := a.docker.Close()
	if err != nil {
		logrus.Errorf("failed to close docker client: %v", err)
	}
}
