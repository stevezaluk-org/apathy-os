package service

import (
	"os"
	"os/exec"

	"github.com/containerd/cgroups/v2/cgroup2"
	"github.com/goccy/go-yaml"
)

// Service Defines a traditional (systemd-style) service that needs to be started. All services are contained within
// cgroups under the "serviced.slice" parent slice. Services should not be used for deploying your applications, use
// containers for that.
//
// Services are considered core to the health of the operating system and do not have a configurable restart policy.
// Services will always be retried 3 times with a 20-second backoff, and if it fails to start in this predefined time
// window, then the operating system will exit
type Service struct {
	// Name The name or identifier for a service
	Name string `yaml:"name"`

	// Binary An absolute path to a binary to execute
	Binary string `yaml:"binary"`

	// Args A list of arguments to call the binary with
	Args []string `yaml:"args"`

	// Env Environment variables to pass to the service. Will be merged with OS set env vars
	Env []string `yaml:"env"`

	// proc Contains the "state" of the service and its associating process
	proc *exec.Cmd

	// manager The cgroup2 manager for this service
	manager *cgroup2.Manager
}

// New Initialize a new Service and return a pointer to it. You must provide an absolute path that the service is
// available at. On initialization, the services cgroup will be created, but it will not be started. Use Service.Start
// for that
func New(path string, parent *cgroup2.Manager) (*Service, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	service := new(Service)
	if err := yaml.Unmarshal(data, service); err != nil {
		return nil, err
	}

	manager, err := parent.NewChild(service.Name+".service", new(cgroup2.Resources))
	if err != nil {
		return nil, err
	}

	service.manager = manager
	service.proc = exec.Command(service.Binary, service.Args...)
	service.proc.Env = append(service.proc.Env, service.Env...)

	return service, nil
}
