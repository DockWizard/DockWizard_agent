package api

import (
	"testing"
	"time"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c := &config.Config{}
	a := New("http://localhost:8080", c, nil)

	require.Equal(t, 25*time.Second, a.client.Timeout)
}
