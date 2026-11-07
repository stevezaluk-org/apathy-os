package env

import (
	"fmt"
	"os"

	"github.com/stevezaluk-org/apathy-os/internal/machinery/task"
)

type EnvironmentVariable struct {
	// Key The name of the environment variable
	Key string `yaml:"key" json:"key"`

	// Value Values of the environment variable
	Value string `yaml:"value" json:"value"`
}

// PriorityGate Always returns 0 on an environment variable since there is no reason that one set of env vars should be
// created before or after another set
func (envVar *EnvironmentVariable) PriorityGate() int {
	return 0
}

// Set Sets an environment variable and returns a wrapped error from it
func Set(env *EnvironmentVariable) error {
	if err := os.Setenv(env.Key, env.Value); err != nil {
		return fmt.Errorf("failed to set environment variable: %v", err)
	}

	return nil
}

// SetAll Sets all environment variables defined. Priority is ignored here.
func SetAll(env []*EnvironmentVariable) error {
	if err := task.New(env).Wait(Set); err != nil {
		return err
	}

	return nil
}
