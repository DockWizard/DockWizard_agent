package dockerstats

import (
	"encoding/json"
	"time"
)

type DockerStats struct {
	Read        time.Time   `json:"read,omitempty"`
	PidsStats   PidsStats   `json:"pids_stats,omitempty"`
	Networks    Networks    `json:"networks,omitempty"`
	MemoryStats MemoryStats `json:"memory_stats,omitempty"`
	BlkioStats  BlkioStats  `json:"blkio_stats,omitempty"`
	CPUStats    CPUStats    `json:"cpu_stats,omitempty"`
	PrecpuStats PrecpuStats `json:"precpu_stats,omitempty"`
}
type PidsStats struct {
	Current int `json:"current,omitempty"`
}
type Eth0 struct {
	RxBytes   int `json:"rx_bytes,omitempty"`
	RxDropped int `json:"rx_dropped,omitempty"`
	RxErrors  int `json:"rx_errors,omitempty"`
	RxPackets int `json:"rx_packets,omitempty"`
	TxBytes   int `json:"tx_bytes,omitempty"`
	TxDropped int `json:"tx_dropped,omitempty"`
	TxErrors  int `json:"tx_errors,omitempty"`
	TxPackets int `json:"tx_packets,omitempty"`
}
type Eth5 struct {
	RxBytes   int `json:"rx_bytes,omitempty"`
	RxDropped int `json:"rx_dropped,omitempty"`
	RxErrors  int `json:"rx_errors,omitempty"`
	RxPackets int `json:"rx_packets,omitempty"`
	TxBytes   int `json:"tx_bytes,omitempty"`
	TxDropped int `json:"tx_dropped,omitempty"`
	TxErrors  int `json:"tx_errors,omitempty"`
	TxPackets int `json:"tx_packets,omitempty"`
}
type Networks struct {
	Eth0 Eth0 `json:"eth0,omitempty"`
	Eth5 Eth5 `json:"eth5,omitempty"`
}
type Stats struct {
	TotalPgmajfault         int `json:"total_pgmajfault,omitempty"`
	Cache                   int `json:"cache,omitempty"`
	MappedFile              int `json:"mapped_file,omitempty"`
	TotalInactiveFile       int `json:"total_inactive_file,omitempty"`
	Pgpgout                 int `json:"pgpgout,omitempty"`
	Rss                     int `json:"rss,omitempty"`
	TotalMappedFile         int `json:"total_mapped_file,omitempty"`
	Writeback               int `json:"writeback,omitempty"`
	Unevictable             int `json:"unevictable,omitempty"`
	Pgpgin                  int `json:"pgpgin,omitempty"`
	TotalUnevictable        int `json:"total_unevictable,omitempty"`
	Pgmajfault              int `json:"pgmajfault,omitempty"`
	TotalRss                int `json:"total_rss,omitempty"`
	TotalRssHuge            int `json:"total_rss_huge,omitempty"`
	TotalWriteback          int `json:"total_writeback,omitempty"`
	TotalInactiveAnon       int `json:"total_inactive_anon,omitempty"`
	RssHuge                 int `json:"rss_huge,omitempty"`
	HierarchicalMemoryLimit int `json:"hierarchical_memory_limit,omitempty"`
	TotalPgfault            int `json:"total_pgfault,omitempty"`
	TotalActiveFile         int `json:"total_active_file,omitempty"`
	ActiveAnon              int `json:"active_anon,omitempty"`
	TotalActiveAnon         int `json:"total_active_anon,omitempty"`
	TotalPgpgout            int `json:"total_pgpgout,omitempty"`
	TotalCache              int `json:"total_cache,omitempty"`
	InactiveAnon            int `json:"inactive_anon,omitempty"`
	ActiveFile              int `json:"active_file,omitempty"`
	Pgfault                 int `json:"pgfault,omitempty"`
	InactiveFile            int `json:"inactive_file,omitempty"`
	TotalPgpgin             int `json:"total_pgpgin,omitempty"`
}
type MemoryStats struct {
	Stats    Stats `json:"stats,omitempty"`
	MaxUsage int   `json:"max_usage,omitempty"`
	Usage    int   `json:"usage,omitempty"`
	Failcnt  int   `json:"failcnt,omitempty"`
	Limit    int   `json:"limit,omitempty"`
}
type BlkioStats struct {
}
type CPUUsage struct {
	PercpuUsage       []int `json:"percpu_usage,omitempty"`
	UsageInUsermode   int   `json:"usage_in_usermode,omitempty"`
	TotalUsage        int   `json:"total_usage,omitempty"`
	UsageInKernelmode int   `json:"usage_in_kernelmode,omitempty"`
}
type ThrottlingData struct {
	Periods          int `json:"periods,omitempty"`
	ThrottledPeriods int `json:"throttled_periods,omitempty"`
	ThrottledTime    int `json:"throttled_time,omitempty"`
}
type CPUStats struct {
	CPUUsage       CPUUsage       `json:"cpu_usage,omitempty"`
	SystemCPUUsage int64          `json:"system_cpu_usage,omitempty"`
	OnlineCpus     int            `json:"online_cpus,omitempty"`
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}
type PrecpuStats struct {
	CPUUsage       CPUUsage       `json:"cpu_usage,omitempty"`
	SystemCPUUsage int64          `json:"system_cpu_usage,omitempty"`
	OnlineCpus     int            `json:"online_cpus,omitempty"`
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}

func Unmarshal(data []byte) (*DockerStats, error) {
	var stats DockerStats
	err := json.Unmarshal(data, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// used_memory = memory_stats.usage - memory_stats.stats.cache
// available_memory = memory_stats.limit
// Memory usage % = (used_memory / available_memory) * 100.0
// cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
// system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
// number_cpus = lenght(cpu_stats.cpu_usage.percpu_usage) or cpu_stats.online_cpus
// CPU usage % = (cpu_delta / system_cpu_delta) * number_cpus * 100.0

func (d *DockerStats) UsedMemory() int {
	return d.MemoryStats.Usage - d.MemoryStats.Stats.Cache
}

func (d *DockerStats) AvailableMemory() int {
	return d.MemoryStats.Limit
}

func (d *DockerStats) MemoryUsagePercentage() float64 {
	return float64(d.UsedMemory()) / float64(d.AvailableMemory()) * 100.0
}

func (d *DockerStats) CpuDelta() int {
	return d.CPUStats.CPUUsage.TotalUsage - d.PrecpuStats.CPUUsage.TotalUsage
}

func (d *DockerStats) SystemCpuDelta() int64 {
	return d.CPUStats.SystemCPUUsage - d.PrecpuStats.SystemCPUUsage
}

func (d *DockerStats) NumberCpus() int {
	return d.CPUStats.OnlineCpus
}

func (d *DockerStats) CpuUsagePercentage() float64 {
	return (float64(d.CpuDelta()) / float64(d.SystemCpuDelta())) * float64(d.NumberCpus()) * 100.0
}
