package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dockwizard/dockwizard/agent/pkg/config"
	"github.com/dockwizard/dockwizard/agent/pkg/data"
	"io"
	"log"
	"net/http"
	"time"
)

type API struct {
	config   *config.Config
	client   *http.Client
	endpoint string
}

type AgentMetadata struct {
	ContainerID   string `json:"container_id"`
	ContainerName string `json:"container_name"`
}

type AgentData struct {
	CPU              int `json:"cpu"`
	MemoryPercentage int `json:"memory_perc"`
	MemoryTotal      int `json:"memory_tot"`
	NetIO            int `json:"net_io"`
	BlockIO          int `json:"block_io"`
}

type AgentObject struct {
	Timestamp time.Time      `json:"timestamp"`
	Metadata  *AgentMetadata `json:"metadata"`
	Data      *AgentData     `json:"data"`
}

type AgentObjectList struct {
	AgentID string         `json:"agent_id"`
	Data    []*AgentObject `json:"data"`
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
				ContainerID:   container.ID,
				ContainerName: container.Name,
			},
			Data: &AgentData{
				CPU:              int(container.CPUUsage),
				MemoryPercentage: container.MemoryUsage,
				MemoryTotal:      0,
				NetIO:            0,
				BlockIO:          0,
			},
		})
	}

	agentObjectList := &AgentObjectList{
		AgentID: a.config.AgentID,
		Data:    list,
	}

	jsonData, err := json.Marshal(agentObjectList)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", a.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
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

	log.Printf("response: %s", string(bts))
	return nil
}
