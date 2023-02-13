package data

type ContainerMetrics struct {
	// ID is the container ID
	ID string `json:"id"`

	// Name is the container name
	Name string `json:"name"`

	// Image is the container image
	Image string `json:"image"`

	// CPUUsage is the CPU usage in percentage
	CPUUsage float64 `json:"cpu_usage"`

	// MemoryUsage is the memory usage in MB
	MemoryUsage int `json:"memory_usage"`

	// MemoryUsagePercentage is the memory usage in percentage
	MemoryUsagePercentage float64 `json:"memory_usage_percentage"`

	// State is the container state
	State string `json:"state"`

	// NetworkIORead is the network IO read in bytes
	NetworkIORead int `json:"network_io_read"`

	// NetworkIOWrite is the network IO write in bytes
	NetworkIOWrite int `json:"network_io_write"`

	// BlockIORead is the block IO read in bytes
	BlockIORead int `json:"block_io_read"`

	// BlockIOWrite is the block IO write in bytes
	BlockIOWrite int `json:"block_io_write"`
}

type Metrics struct {
	Container []*ContainerMetrics
}
