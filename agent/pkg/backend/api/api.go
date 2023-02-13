package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dockwizard/dockwizard/agent/pkg/config"
	"github.com/dockwizard/dockwizard/agent/pkg/data"
)

type API struct {
	config   *config.Config
	client   *http.Client
	endpoint string
}

type AgentMetadata struct {
	ContainerID    string `json:"container_id"`
	ContainerName  string `json:"container_name"`
	ContainerImage string `json:"container_image"`
}

type AgentData struct {
	CPU              float64 `json:"cpu"`
	MemoryPercentage float64 `json:"memory_perc"`
	MemoryTotal      int     `json:"memory_tot"`
	TotalRx          int     `json:"total_rx"`
	TotalTx          int     `json:"total_tx"`
	IoRead           int     `json:"io_read"`
	IoWrite          int     `json:"io_write"`
}

type AgentObject struct {
	Timestamp time.Time      `json:"timestamp"`
	Metadata  *AgentMetadata `json:"metadata"`
	Data      *AgentData     `json:"data"`
}

type AgentObjectList struct {
	Data []*AgentObject `json:"data"`
}

func New(endpoint string, config *config.Config) *API {
	client := &http.Client{
		Timeout: 25 * time.Second,
	}
	return &API{
		config:   config,
		client:   client,
		endpoint: endpoint,
	}
}

func (a *API) SendData(metrics *data.Metrics) error {
	var list []*AgentObject
	for _, container := range metrics.Container {
		list = append(list, &AgentObject{
			Timestamp: time.Now(),
			Metadata: &AgentMetadata{
				ContainerID:    container.ID,
				ContainerName:  container.Name,
				ContainerImage: container.Image,
			},
			Data: &AgentData{
				CPU:              container.CPUUsage,
				MemoryTotal:      container.MemoryUsage,
				MemoryPercentage: container.MemoryUsagePercentage,
				TotalRx:          container.NetworkIORead,
				TotalTx:          container.NetworkIOWrite,
				IoRead:           container.BlockIORead,
				IoWrite:          container.BlockIOWrite,
			},
		})
	}

	agentObjectList := &AgentObjectList{
		Data: list,
	}

	jsonData, err := json.Marshal(agentObjectList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.config.APIKey))
	res, err := a.client.Do(req)
	if err != nil {
		return err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("response: %s", string(bts))
	}

	return nil
}
