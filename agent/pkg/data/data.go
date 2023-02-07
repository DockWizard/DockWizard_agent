package data

type ContainerMetrics struct {
	// ID is the container ID
	ID string `json:"id"`

	// Name is the container name
	Name string `json:"name"`

	// CPUUsage is the CPU usage in percentage
	CPUUsage float64 `json:"cpu_usage"`

	// MemoryUsage is the memory usage in MB
	MemoryUsage int `json:"memory_usage"`

	// State is the container state
	State string `json:"state"`
}

type Metrics struct {
	Container []*ContainerMetrics
}
