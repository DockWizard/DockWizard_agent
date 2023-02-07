package stdout

import (
	"encoding/json"
	"github.com/dockwizard/dockwizard/agent/pkg/data"
	"os"
)

type Stdout struct{}

func (s *Stdout) SendData(metrics *data.Metrics) error {
	return json.NewEncoder(os.Stdout).Encode(metrics)
}
