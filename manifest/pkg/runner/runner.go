package runner

import "github.com/stevezaluk/immutable-os/manifest/pkg/config"

// Runner Handles actually processing a config file and executing its required steps
type Runner struct {
	config *config.Config
}

// New Initializes a new Runner and returns a pointer to it
func New(config *config.Config) *Runner {
	return &Runner{config: config}
}
