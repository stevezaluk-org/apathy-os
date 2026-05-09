package service

import (
	"context"
	"log"
	"os"
	"os/exec"
	"syscall"
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

	// ctx Context controlling concurrency for service starting/stopping
	ctx context.Context

	// cancel Function for stopping executing
	cancel context.CancelFunc
}

func (service *Service) Stop(ctx context.Context) error {
	if service.proc.Process == nil {
		return nil // Service has not been started exit early
	}

	log.Printf("Sending SIGTERM to service %s with PID %d", service.Name, service.proc.Process.Pid)
	if err := service.proc.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Service failed to respond to SIGTERM. Service will be SIGKILL'ed and cgroup frozen")
			if err := service.manager.Freeze(); err != nil {
				return err
			}

			return nil
		default:
			if err := service.proc.ProcessState; err != nil {
				log.Printf("Service %s stopped. Exit code: %d", service.Name, service.proc.ProcessState.ExitCode())
				return nil
			}
		}
	}
}

// Start Populates the service.proc field with an exec.Cmd pointer and starts the service. Does not automatically transition
// to a Watch loop with automatic restarts. The context passed to this function should be derived from the top level
// parent context so that you can stop services from trying to restart when system shutdown is received
func (service *Service) Start(ctx context.Context) error {
	service.proc = exec.CommandContext(ctx, service.Binary, service.Args...)
	service.proc.Cancel = func() error {
		cleanUpCtx, cleanUpCancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cleanUpCancel()

		if err := service.Stop(cleanUpCtx); err != nil {
			return err
		}

		return nil
	}

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

	// migrate the service to its cgroup
	if err := service.manager.AddProc(uint64(service.proc.Process.Pid)); err != nil {
		return err
	}

	return nil
}

// Watch Initiates a Watch loop on the specified service and starts it. Provides an automatic restart loop for the service.
// A new child context is  derived from the parent context that is passed to this function in order to ensure that
// services do not attempt to restart after system shutdown has been initiated
// TODO: Service accepts input (when os.Stdin) is mapped even after Service.Stop is called
func (service *Service) Watch(ctx context.Context) {
	child, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := range 3 {
		select {
		case <-child.Done():
			log.Println("Service watch cancelled")
			return
		default:
			log.Printf("Attempting to start service %s on attempt %d", service.Name, i+1)
			if err := service.Start(child); err != nil {
				log.Printf("Failed to start service %s: %v", service.Name, err)
			}

			log.Printf("Service %s started with PID %d", service.Name, service.proc.Process.Pid)
			if err := service.proc.Wait(); err != nil {
				log.Printf("Service %s died with error: %v", service.Name, err)
			}
		}

		time.Sleep(20 * time.Second) // back off
	}
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
