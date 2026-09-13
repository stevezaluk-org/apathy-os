package manifest

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/env"
	"github.com/immutable-os/immutable-os/src/go-mount/pkg/mount"
	"github.com/immutable-os/immutable-os/src/go-mount/pkg/overlay"
)

// Manifest An abstraction of a manifest file
// TODO: Needs to support merging config files
type Manifest struct {
	// Mounts The `mounts` section of the manifest file
	Mounts []*mount.Mount `yaml:"mounts" json:"mounts"`

	// Envs EnvironmentVariables that should be created
	Envs []*env.EnvironmentVariable `yaml:"env" json:"env"`

	// Overlays Overlay mounts that should be created
	Overlays []*overlay.Overlay `yaml:"overlays" json:"overlays"`
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
