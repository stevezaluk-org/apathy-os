package service

// ProcessService Defines a basic binary that should be treated and managed as if it was a system service
type ProcessService struct {
	// Binary An absolute unix path to the binary that should be executed
	Binary string `yaml:"bin"`

	// Args CLI arguments to pass to the binary
	Args []string `yaml:"args"`

	// Exec If set to true, unix.Exec is executed to start this binary
	Exec bool `yaml:"exec"`

	// Priority The priority gate for the process
	Priority int `yaml:"priority"`
}

// PriorityGate The priority gate for the process
func (processService *ProcessService) PriorityGate() int {
	return processService.Priority
}

// Services Defines the services config section
type Services struct {
	ProcessServices []*ProcessService `yaml:"processes"`
}
