package agent

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/dockwizard/dockwizard_agent/agent/internal/testutils"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/backend/stdout"
	"github.com/dockwizard/dockwizard_agent/agent/pkg/config"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func newAgent(m client.APIClient) *Agent {
	return New(&config.Config{}, stdout.New(), m)
}

func TestGetDockerContainerMetricsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := testutils.NewMockAPIClient(ctrl)
	agent := newAgent(m)

	m.
		EXPECT().
		ContainerList(gomock.Any(), types.ContainerListOptions{}).
		Return([]types.Container{}, nil)

	metrics, err := agent.getDockerContainerMetrics()
	require.Nil(t, err)
	require.Equal(t, 0, len(metrics))
}

func TestGetDockerContainerMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)

	m := testutils.NewMockAPIClient(ctrl)
	agent := newAgent(m)

	containers := []types.Container{
		{
			ID: "1",
			Names: []string{
				"test",
			},
		},
	}

	var b bytes.Buffer
	b.Write([]byte(`
	{
		"read": "2023-02-20T10:03:01.998224131Z",
		"preread": "0001-01-01T00:00:00Z",
		"pids_stats": {
			"current": 1,
			"limit": 18446744073709552000
		},
		"blkio_stats": {
			"io_service_bytes_recursive": [
				{
					"major": 254,
					"minor": 0,
					"op": "read",
					"value": 3383296
				},
				{
					"major": 254,
					"minor": 0,
					"op": "write",
					"value": 0
				}
			],
			"io_serviced_recursive": null,
			"io_queue_recursive": null,
			"io_service_time_recursive": null,
			"io_wait_time_recursive": null,
			"io_merged_recursive": null,
			"io_time_recursive": null,
			"sectors_recursive": null
		},
		"num_procs": 0,
		"storage_stats": {},
		"cpu_stats": {
			"cpu_usage": {
				"total_usage": 168791000,
				"usage_in_kernelmode": 98048000,
				"usage_in_usermode": 70743000
			},
			"system_cpu_usage": 2006893600000000,
			"online_cpus": 4,
			"throttling_data": {
				"periods": 0,
				"throttled_periods": 0,
				"throttled_time": 0
			}
		},
		"precpu_stats": {
			"cpu_usage": {
				"total_usage": 0,
				"usage_in_kernelmode": 0,
				"usage_in_usermode": 0
			},
			"throttling_data": {
				"periods": 0,
				"throttled_periods": 0,
				"throttled_time": 0
			}
		},
		"memory_stats": {
			"usage": 4194304,
			"stats": {
				"active_anon": 0,
				"active_file": 1892352,
				"anon": 811008,
				"anon_thp": 0,
				"file": 3108864,
				"file_dirty": 0,
				"file_mapped": 1892352,
				"file_writeback": 0,
				"inactive_anon": 811008,
				"inactive_file": 1253376,
				"kernel_stack": 0,
				"pgactivate": 330,
				"pgdeactivate": 0,
				"pgfault": 10626,
				"pglazyfree": 0,
				"pglazyfreed": 0,
				"pgmajfault": 0,
				"pgrefill": 0,
				"pgscan": 421,
				"pgsteal": 0,
				"shmem": 0,
				"slab": 262440,
				"slab_reclaimable": 184,
				"slab_unreclaimable": 262256,
				"sock": 0,
				"thp_collapse_alloc": 0,
				"thp_fault_alloc": 0,
				"unevictable": 0,
				"workingset_activate": 0,
				"workingset_nodereclaim": 0,
				"workingset_refault": 0
			},
			"limit": 12544401408
		},
		"name": "/test",
		"id": "1",
		"networks": {
			"eth0": {
				"rx_bytes": 37188,
				"rx_packets": 500,
				"rx_errors": 0,
				"rx_dropped": 0,
				"tx_bytes": 10036,
				"tx_packets": 142,
				"tx_errors": 0,
				"tx_dropped": 0
			}
		}
	}
`))
	containerStats := types.ContainerStats{
		OSType: "linux",
		Body:   io.NopCloser(&b),
	}

	m.
		EXPECT().
		ContainerList(gomock.Any(), types.ContainerListOptions{}).
		Return(containers, nil)

	m.
		EXPECT().
		ContainerStatsOneShot(gomock.Any(), "1").
		Return(containerStats, nil)

	metrics, err := agent.getDockerContainerMetrics()
	require.Nil(t, err)
	require.Equal(t, 1, len(metrics))
	require.Equal(t, "1", metrics[0].ID)
	require.Equal(t, "test", metrics[0].Name)
}

func TestSleep(t *testing.T) {
	c := &config.Config{
		UpdateFrequency: 1,
	}
	agent := New(c, stdout.New(), nil)

	timeNow := time.Now()
	agent.sleep()
	timeAfter := time.Now()

	require.True(t, timeAfter.Sub(timeNow) >= time.Second)
}
