package service

import (
	"context"
	"log"
	"os"
	"os/exec"
	"time"

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

	// MapConsole If set to true, map the processes Stdout/Stdin/Stderr to there os equivalents. Allows you to actually
	// interact with a service that wants input from stdin
	MapConsole bool `yaml:"mapConsole"`

	// proc Contains the "state" of the service and its associating process
	proc *exec.Cmd

	// manager The cgroup2 manager for this service
	manager *cgroup2.Manager
}

// Watch Starts the core wait/restart loop for the service. Services do not have a configurable restart policy as they
// are treated as core to the operating system. Services will try and restart a maximum of 3 times before Watch gives
// up and trying to restart them. Backoff in between restarts is set to 5 seconds
func (service *Service) Watch(ctx context.Context, cancel context.CancelCauseFunc) {
	log.Printf("Watching service %s", service.Name)

	for i := range 3 {
		select {
		case <-ctx.Done():
			cancel(ctx.Err()) // context was canceled, ensure that we have proper clean up here
		default:
			if err := service.proc.Wait(); err != nil {
				log.Printf("Service %s (%d) died with error: %s", service.Name, service.proc.Process.Pid, err)
			}

			time.Sleep(5 * time.Second) // backoff

			// TODO: This does not free the process from its cgroup, which should be changed
			log.Printf("Restarting service %s. Try %d", service.Name, i)
			if err := service.Start(); err != nil {
				log.Printf("Got error trying to restart service %s on try %d: %s", service.Name, i, err)
				if i == 3 {
					cancel(err) // we hit our limit on restarting the service. proper clean up here
				}
			}
		}
	}
}

// Start Starts the specified service, and provided that it does not immediately crash, its process is migrated to the
// services cgroup.
func (service *Service) Start() error {
	service.proc = exec.Command(service.Binary, service.Args...)
	service.proc.Env = append(service.proc.Env, service.Env...)

	// This is kind of a hacky way of doing this and we really should be using something like containerd/console for
	// creating a proper virtual console. To be replaced in the future
	// TODO: Replace with github.com/containerd/console
	if service.MapConsole {
		service.proc.Stdout = os.Stdout
		service.proc.Stdin = os.Stdin
		service.proc.Stderr = os.Stderr
	}

	if err := service.proc.Start(); err != nil {
		return err
	}

	if service.proc.Process == nil {
		return nil // service failed to start here. Actually process the error
	}

	// migrate the service to its cgroup
	if err := service.manager.AddProc(uint64(service.proc.Process.Pid)); err != nil {
		return err
	}

	log.Printf("Service %s started with PID %d", service.Name, service.proc.Process.Pid)

	return nil
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

	return service, nil
}
