package service

// Service Defines a traditional (systemd-style) service that needs to be started. All services are contained within
// cgroups under the "serviced.slice" parent slice. Services should not be used for deploying your applications, use
// containers for that.
//
// Services are considered core to the health of the operating system and do not have a configurable
// restart policy. Services will always be retried 3 times with a 20-second backoff, and if it fails to start in this
// predefined time window, then the operating system will exit
type Service struct {
	// Name The name or identifier for a service
	Name string `yaml:"name"`

	// Binary An absolute path to a binary to execute
	Binary string `yaml:"binary"`

	// Args A list of arguments to call the binary with
	Args []string `yaml:"args"`

	// Env Environment variables to pass to the service. Will be merged with OS set env vars
	Env []string `yaml:"env"`
}
