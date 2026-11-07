package runner

import (
	"github.com/stevezaluk-org/apathy-os/internal/manifest"
	"github.com/stevezaluk-org/apathy-os/internal/manifest/env"
	"github.com/stevezaluk-org/apathy-os/internal/mount/filesystem"
	"github.com/stevezaluk-org/apathy-os/internal/mount/overlay"
)

// Runner Handles actually processing a config file and executing its required steps
type Runner struct {
	// config Caller passed config.Config struct.
	config *manifest.Manifest
}

// Go Start the runner and process the manifest config
func (runner *Runner) Go() error {
	if err := env.SetAll(runner.config.Envs); err != nil {
		return err
	}

	if err := filesystem.MountAll(runner.config.Mounts); err != nil {
		return err
	}

	if err := overlay.MountAll(runner.config.Overlays); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Runner and returns a pointer to it
func New(config *manifest.Manifest) *Runner {
	return &Runner{
		config: config,
	}
}
