package dockerstats_test

import (
	"testing"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/dockerstats"
	"github.com/stretchr/testify/require"
)

var testPayload = `
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
`

func TestUnmarshal(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 1, stats.PidsStats.Current)
	require.Equal(t, 4194304, stats.MemoryStats.Usage)
}

func TestUsedMemory(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 4194304, stats.UsedMemory())
}

func TestAvailableMemory(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 12544401408, stats.AvailableMemory())
}

func TestMemoryUsagePercentage(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 0.03343566475260547, stats.MemoryUsagePercentage())
}

func TestCpuDelta(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 168791000, stats.CpuDelta())
}

func TestSystemCpuDelta(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, int64(2006893600000000), stats.SystemCpuDelta())
}

func TestNumberCpus(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 4, stats.NumberCpus())
}

func TestCpuUsagePercentage(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	require.Equal(t, 3.364224192054825e-05, stats.CpuUsagePercentage())
}

func TestNetworkStats(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	rx, tx := stats.NetworkStats()
	require.Equal(t, 37188.0, rx)
	require.Equal(t, 10036.0, tx)
}

func TestDiskStats(t *testing.T) {
	stats, err := dockerstats.Unmarshal([]byte(testPayload))
	require.Nil(t, err)

	read, write := stats.DiskStats()
	require.Equal(t, uint64(0x33a000), read)
	require.Equal(t, uint64(0x0), write)
}
