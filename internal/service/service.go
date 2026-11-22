package service

import (
	"encoding/json"
	"fmt"
	"os"
)

// Service Represents a service that should be started
type Service struct {
	// Name The name of the service
	Name string `json:"name"`

	// Exec The binary to exec as a part of starting the service
	Exec string `json:"exec"`

	// Arguments A list of arguments to append to Service.Exec when starting the service
	Args []string `json:"args"`

	// Env Environment variables that the service should start with
	Env []string `json:"env"`

	// Priority A number dictating when the service should start, relative to other services
	Priority int `json:"priority"`

	// AllocateTty If set to true, provide the service a TTY and map it to stdin
	AllocateTty bool `json:"allocate_tty"`

	// MaxRestartCount The maximum number of times that this service will be autoamtically restarted on crash
	MaxRestartCount int `json:"max_restart_count"`
}

// String Returns the service represented as a string
func (service *Service) String() string {
	return fmt.Sprintf("%s (%d) -> ./%s %s",
		service.Name, service.Priority, service.Exec, service.Args)
}

// PriorityGate Returns the priority of the service. Lower number means the service starts sooner, higher means later,
// and services with the same number are started concurrently
func (service *Service) PriorityGate() int {
	return service.Priority
}

// FromFile Reads a JSON file at its respective path and unmarshals it to a service
func FromFile(path string) (*Service, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read service file: %v", err)
	}

	ret := new(Service)
	if err := json.Unmarshal(contents, ret); err != nil {
		return nil, fmt.Errorf("unmarshal error when reading service file: %v", err)
	}

	return ret, nil
}
