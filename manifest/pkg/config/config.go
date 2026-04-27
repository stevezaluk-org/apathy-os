package config

import (
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// Config An abstraction of a manifest file
// TODO: Needs to support merging config files
type Config struct {
	// Mounts The `mounts` section of the manifest file
	Mounts MountConfig `yaml:"mounts"`

	// Envs EnvironmentVariables that should be created
	Envs []EnvironmentVariable `yaml:"envs"`
}

// New Initialize a Config struct and return a pointer to it
func New() *Config {
	return &Config{
		Mounts: make(MountConfig, 0),
	}
}

// Read Reads a config file at the specified Unix path, unmarshalls it into Config, and returns a pointer to it
func Read(path string) (*Config, error) {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}

	config := New()

	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(contents, config); err != nil {
		return nil, err
	}

	return config, err
}
