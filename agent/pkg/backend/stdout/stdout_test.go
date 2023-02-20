package stdout

import (
	"os"
	"testing"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/data"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	s := New()
	require.Equal(t, os.Stdout, s.stdout)
}

func TestSendData(t *testing.T) {
	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)

	s := &stdout{
		stdout: tmp,
	}
	s.SendData(&data.Metrics{
		Container: []*data.ContainerMetrics{},
	})

	_, err = tmp.Seek(0, 0)
	require.Nil(t, err)

	b := make([]byte, 16)
	_, err = tmp.Read(b)
	require.Nil(t, err)

	require.Equal(t, `{"Container":[]}`, string(b))
}
