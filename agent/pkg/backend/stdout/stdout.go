package stdout

import (
	"encoding/json"
	"os"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/data"
)

type Stdout struct{}

func (s *Stdout) SendData(metrics *data.Metrics) error {
	return json.NewEncoder(os.Stdout).Encode(metrics)
}
