package agent

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dockwizard/dockwizard/agent/pkg/backend"
	"github.com/dockwizard/dockwizard/agent/pkg/config"
	"github.com/dockwizard/dockwizard/agent/pkg/data"
	"github.com/dockwizard/dockwizard/agent/pkg/dockerstats"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

type Agent struct {
	Config *config.Config

	ctx     context.Context
	cancel  func()
	log     *logrus.Logger
	buffer  []*data.Metrics
	docker  *client.Client
	backend backend.Backend
}

func New(c *config.Config, b backend.Backend) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		cancel()
		return nil, err
	}

	return &Agent{
		Config:  c,
		ctx:     ctx,
		cancel:  cancel,
		log:     logrus.New(),
		buffer:  []*data.Metrics{},
		docker:  cli,
		backend: b,
	}, nil
}

func (a *Agent) getDockerContainerMetrics() ([]*data.ContainerMetrics, error) {
	var ret []*data.ContainerMetrics

	// Get the metrics
	allContainers, err := a.docker.ContainerList(a.ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	// Get the metrics for each container
	for _, container := range allContainers {
		stats, err := a.docker.ContainerStatsOneShot(a.ctx, container.ID)
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

		ret = append(ret, &data.ContainerMetrics{
			ID:          container.ID,
			Name:        strings.TrimPrefix(container.Names[0], "/"),
			CPUUsage:    parsedStats.CpuUsagePercentage(),
			MemoryUsage: parsedStats.UsedMemory(),
			State:       container.State,
		})
	}

	return ret, nil
}

func (a *Agent) sleep() {
	time.Sleep(time.Duration(a.Config.PollInterval) * time.Second)
}

func (a *Agent) Run() {
	for {
		// Sleep for the poll interval
		a.sleep()

		// If canceled then break
		if a.ctx.Err() == context.Canceled {
			break
		}

		containerMetrics, err := a.getDockerContainerMetrics()
		if err != nil {
			a.log.Errorf("error getting container metrics: %v", err)
			continue
		}

		// Send the metrics to the backend
		metrics := &data.Metrics{
			Container: containerMetrics,
		}
		err = a.backend.SendData(metrics)
		if err != nil {
			a.log.Errorf("could not send metrics to backend: %v", err)
		}
	}

	err := a.docker.Close()
	if err != nil {
		logrus.Errorf("failed to close docker client: %v", err)
	}
}
