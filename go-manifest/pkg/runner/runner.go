package runner

import (
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/env"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/manifest"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/mount"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/service"
)

// Runner Handles actually processing a config file and executing its required steps
type Runner struct {
	// config Caller passed config.Config struct.
	config *manifest.Manifest

	// MountHandler Handles processing the `mount` section of a manifest config
	MountHandler *mount.Handler

	// EnvHandler Handles processing the `env` section of a manifest config
	EnvHandler *env.Handler

	// ServiceHandler Handles starting and supervising system processes
	ServiceHandler *service.Handler
}

// Go Start the runner and process the manifest config
func (runner *Runner) Go() error {
	if err := runner.EnvHandler.SetAll(); err != nil {
		return err
	}

	if err := runner.MountHandler.MountAll(false); err != nil {
		return err
	}

	if err := runner.ServiceHandler.StartAll(); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Runner and returns a pointer to it
func New(config *manifest.Manifest) *Runner {
	return &Runner{
		config:         config,
		MountHandler:   mount.NewHandler(config.Mounts),
		EnvHandler:     env.NewHandler(config.Envs),
		ServiceHandler: service.NewHandler(config.Services.ProcessServices),
	}
}
