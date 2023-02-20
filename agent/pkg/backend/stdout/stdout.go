package stdout

import (
	"encoding/json"
	"os"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/data"
)

type stdout struct {
	stdout *os.File
}

func New() *stdout {
	return &stdout{
		stdout: os.Stdout,
	}
}

func (s *stdout) SendData(metrics *data.Metrics) error {
	return json.NewEncoder(s.stdout).Encode(metrics)
}
