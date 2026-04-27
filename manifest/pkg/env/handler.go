package env

import (
	"errors"
	"os"

	internalErrors "github.com/stevezaluk/immutable-os/manifest/internal/errors"
	"github.com/stevezaluk/immutable-os/manifest/pkg/config"
	"golang.org/x/sync/errgroup"
)

var ErrSetEnvFailed = errors.New("set environment variable failed")

// Handler Processes the `env` section of a config file. Concurrently sets system wide environment variables
type Handler struct {
	envs []*config.EnvironmentVariable
}

// SetAll Set's all environment variables defined in Handler.envs
func (handler *Handler) SetAll() error {
	wg := new(errgroup.Group)

	for _, env := range handler.envs {
		wg.Go(func() error {
			if err := os.Setenv(env.Key, env.Value); err != nil {
				return internalErrors.Wrap(ErrSetEnvFailed, err)
			}

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return nil
}

// NewHandler Initializes a new Handler and returns a pointer to it
func NewHandler(envs []*config.EnvironmentVariable) *Handler {
	return &Handler{
		envs: envs,
	}
}
