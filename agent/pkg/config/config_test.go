package config_test

import (
	"os"
	"testing"

	"github.com/dockwizard/dockwizard_agent/agent/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	in := `---
backend: stdout
update_frequency: 2
containers: []
`
	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	c, err := config.Read(tmp.Name())
	require.Nil(t, err)
	require.Equal(t, "stdout", c.Backend)
	require.Equal(t, 2, c.UpdateFrequency)
	require.Equal(t, 0, len(c.Containers))
}

func TestReadInvalid(t *testing.T) {
	in := `---
x: x: 1
`
	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	_, err = config.Read(tmp.Name())
	require.Errorf(t, err, "yaml: line 2: mapping values are not allowed in this context")
}

func TestReadNotExists(t *testing.T) {
	_, err := config.Read("not-exists")
	require.Errorf(t, err, "stat not-exists: no such file or directory")
}

func TestReadNoBackend(t *testing.T) {
	in := `---
update_frequency: 2
containers: []
`

	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	_, err = config.Read(tmp.Name())
	require.Errorf(t, err, "backend is required, use stdout to print to stdout")
}

func TestReadInvalidUpdateFrequency(t *testing.T) {
	in := `---
backend: stdout
update_frequency: -1
containers: []
`

	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	_, err = config.Read(tmp.Name())
	require.Errorf(t, err, "update frequency must be at least 2 seconds")
}

func TestReadAPIEndpoin(t *testing.T) {
	in := `---
backend: api
api_endpoint: http://localhost:8080
api_key: 123
update_frequency: 2
containers: []
`

	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	_, err = config.Read(tmp.Name())
	require.Nil(t, err)
}

func TestReadAPIEndpointInvalid(t *testing.T) {
	in := `---
backend: api
update_frequency: 2
containers: []
`

	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	_, err = config.Read(tmp.Name())
	require.Errorf(t, err, "api endpoint is required when backend is api")
}

func TestReadAPIKeyInvalid(t *testing.T) {
	in := `---
backend: api
api_endpoint: http://localhost:8080
update_frequency: 2
containers: []
`

	tmp, err := os.CreateTemp("", "")
	require.Nil(t, err)
	_, err = tmp.Write([]byte(in))
	require.Nil(t, err)

	_, err = config.Read(tmp.Name())
	require.Errorf(t, err, "api key is required when backend is api")
}
