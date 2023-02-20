package testutils

//go:generate mockgen -destination agent/internal/testutils/dockermock.go -package testutils github.com/docker/docker/client APIClient
