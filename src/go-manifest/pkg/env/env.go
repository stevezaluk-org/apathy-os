package env

type EnvironmentVariable struct {
	// Key The name of the environment variable
	Key string `yaml:"key"`

	// Value Values of the environment variable
	Value string `yaml:"value"`
}

// PriorityGate Always returns 0 on an environment variable since there is no reason that one set of env vars should be
// created before or after another set
func (envVar *EnvironmentVariable) PriorityGate() int {
	return 0
}
