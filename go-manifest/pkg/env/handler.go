package env

import (
	"errors"
	"os"

	"github.com/stevezaluk/immutable-os/go-manifest/pkg/engine"
)

var ErrSetEnvFailed = errors.New("set environment variable failed")

// Handler Processes the `env` section of a config file. Concurrently sets systemwide environment variables
type Handler struct {
	engine *engine.Engine[*EnvironmentVariable]

	// envs Parsed environment variables that should be set systemwide
	envs []*EnvironmentVariable
}

// setEnv Sets an environment variable and returns a wrapped error from it
func (handler *Handler) setEnv(env *EnvironmentVariable) error {
	if err := os.Setenv(env.Key, env.Value); err != nil {
		return err
	}

	return nil
}

// SetAll Set's all environment variables defined in Handler.envs
func (handler *Handler) SetAll() error {
	if err := handler.engine.Wait(handler.setEnv); err != nil {
		return err
	}

	return nil
}

// NewHandler Initializes a new Handler and returns a pointer to it
func NewHandler(envs []*EnvironmentVariable) *Handler {
	return &Handler{
		envs:   envs,
		engine: engine.New(envs),
	}
}
