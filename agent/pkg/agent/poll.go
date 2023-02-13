package agent

import (
	"context"
	"io"
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

		rx, tx := calculateNetwork(parsedStats.Networks)
		read, write := calculateBlockIO(parsedStats.BlkioStats)

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

// From: https://github.com/docker/cli/blob/c1733165159c08101adb0e1f120c7181533550ef/cli/command/container/stats_helpers.go#LL217-L225C2
func calculateNetwork(network map[string]dockerstats.Network) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}

// From: https://github.com/docker/cli/blob/c1733165159c08101adb0e1f120c7181533550ef/cli/command/container/stats_helpers.go#LL201-L215C2
func calculateBlockIO(blkio types.BlkioStats) (uint64, uint64) {
	var blkRead, blkWrite uint64
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		if len(bioEntry.Op) == 0 {
			continue
		}
		switch bioEntry.Op[0] {
		case 'r', 'R':
			blkRead = blkRead + bioEntry.Value
		case 'w', 'W':
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return blkRead, blkWrite
}
