package service

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/containerd/cgroups/v2/cgroup2"
)

const (
	cgroupMount string = "/sys/fs/cgroup"

	// DefaultServicePath The default unix path that services will be discovered within
	DefaultServicePath string = "/usr/share/service.d"
)

// Supervisor An abstraction for managing the lifecycle of services. Provides core centralized functionality for evaluating
// the health of services defined in /etc/service.d/*.yaml.
type Supervisor struct {
	// processList Universal process list for tracking running/stopped services
	processList map[string]*Service

	// parent The parent cgroup that all services are placed within. serviced.slice by default
	parent *cgroup2.Manager

	// mu Mutex lock protecting concurrent writes to the process list
	mu sync.RWMutex
}

func (supervisor *Supervisor) EventLoop(ctx context.Context) {
	for _, service := range supervisor.processList {
		log.Printf("Starting watch loop on %s", service.Name)
		go service.Start(ctx)
	}
}

// NewSupervisor Initializes a new Supervisor and returns a pointer to it. NewSupervisor will initialize the parent
// cgroup in which all processes will be placed within, so you should only call this once. Additionally, NewSupervisor
// walks the Unix path provided in its first argument trying to identify files that end in .yaml. If a .yaml file is found
// it is treated as a service file and a new Service is constructed based on this file.
func NewSupervisor(path string) (*Supervisor, error) {
	// Most immutable-os builds don't define more than 5 services so we are safe to pre-allocate for up to 5 here
	services := make(map[string]*Service, 5)

	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), -1)
	}

	parent, err := cgroup2.NewManager(cgroupMount, "/serviced.slice", new(cgroup2.Resources))
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".yaml") {
			service, err := New(path, parent)
			if err != nil {
				return err
			}

			log.Printf("Service queued for start: %s", path)
			services[service.Name] = service
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	supervisor := &Supervisor{
		processList: services,
		parent:      parent,
		mu:          sync.RWMutex{},
	}

	return supervisor, nil
}
