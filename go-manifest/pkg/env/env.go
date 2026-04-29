package env

type EnvironmentVariable struct {
	// Key The name of the environment variable
	Key string `yaml:"key"`

	// Value Values of the environment variable
	Value string `yaml:"value"`
}
