package backend

import (
	"github.com/dockwizard/dockwizard/agent/pkg/data"
)

type Backend interface {
	SendData(metrics *data.Metrics) error
}
