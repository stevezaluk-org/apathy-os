package runner

import (
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/env"
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/manifest"
	"github.com/immutable-os/immutable-os/src/go-mount/pkg/filesystem"
)

// Runner Handles actually processing a config file and executing its required steps
type Runner struct {
	// config Caller passed config.Config struct.
	config *manifest.Manifest

	// EnvHandler Handles processing the `env` section of a manifest config
	EnvHandler *env.Handler
}

// Go Start the runner and process the manifest config
func (runner *Runner) Go() error {
	if err := runner.EnvHandler.SetAll(); err != nil {
		return err
	}

	if err := filesystem.MountAll(runner.config.Mounts); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Runner and returns a pointer to it
func New(config *manifest.Manifest) *Runner {
	return &Runner{
		config:     config,
		EnvHandler: env.NewHandler(config.Envs),
	}
}
