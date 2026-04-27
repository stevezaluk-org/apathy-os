package runner

import (
	"github.com/stevezaluk/immutable-os/manifest/pkg/config"
	"github.com/stevezaluk/immutable-os/manifest/pkg/env"
	"github.com/stevezaluk/immutable-os/manifest/pkg/mount"
)

// Runner Handles actually processing a config file and executing its required steps
type Runner struct {
	// config Caller passed config.Config struct.
	config *config.Config

	// MountHandler Handles processing the `mount` section of a manifest config
	MountHandler *mount.Handler

	// EnvHandler Handles processing the `env` section of a manifest config
	EnvHandler *env.Handler
}

// Go Start the runner and process the manifest config
func (runner *Runner) Go() error {
	if err := runner.EnvHandler.SetAll(); err != nil {
		return err
	}

	if err := runner.MountHandler.MountAll(false); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Runner and returns a pointer to it
func New(config *config.Config) *Runner {
	return &Runner{
		config:       config,
		MountHandler: mount.NewHandler(config.Mounts),
		EnvHandler:   env.NewHandler(config.Envs),
	}
}
