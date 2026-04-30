package manifest

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/env"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/mount"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/service"
)

// Manifest An abstraction of a manifest file
// TODO: Needs to support merging config files
type Manifest struct {
	// Mounts The `mounts` section of the manifest file
	Mounts []*mount.Mount `yaml:"mounts"`

	// Envs EnvironmentVariables that should be created
	Envs []*env.EnvironmentVariable `yaml:"env"`

	// Services Service's that should be started after the root filesystem is mounted
	Services *service.Services `yaml:"services"`
}

// Read Reads a manifest config file at the specified Unix path, unmarshalls it into Config, and returns a pointer to it
func Read(path string) (*Manifest, error) {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}

	manifestFile := new(Manifest)

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(contents, manifestFile); err != nil {
		return nil, err
	}

	return manifestFile, err
}
