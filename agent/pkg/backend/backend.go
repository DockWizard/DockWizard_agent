package backend

import (
	"github.com/dockwizard/dockwizard_agent/agent/pkg/data"
)

type Backend interface {
	SendData(metrics *data.Metrics) error
}
