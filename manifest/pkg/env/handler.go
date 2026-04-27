package env

import "github.com/stevezaluk/immutable-os/manifest/pkg/config"

// Handler Processes the `env` section of a config file. Concurrently sets system wide environment variables
type Handler struct {
	envs []*config.EnvironmentVariable
}

// NewHandler Initializes a new Handler and returns a pointer to it
func NewHandler(envs []*config.EnvironmentVariable) *Handler {
	return &Handler{
		envs: envs,
	}
}
