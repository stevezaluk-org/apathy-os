package service

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/cgroups/v2/cgroup2"
	"github.com/stevezaluk-org/apathy-os/internal/machinery/task"
)

// TODO: These should be exteralized
const (
	ServicePath      = "/usr/share/service.d"
	CgroupMountPoint = "/sys/fs/cgroup"
)

// Manager Global manager for all services running on the systems. Each manager is associated with one or more services, named by its
// namespace. Each manager gets a cgroup allocated to it, with each service under the managers control getting its own child cgroup.
type Manager struct {
	// namespace The namespace for this manager
	namespace string

	// services A slice of services that the manager has enumerated
	services []*Supervisor

	// engine The task engine responsible for starting services concurrently according to there priority
	engine *task.Engine[*Supervisor]

	// parent The parent cgroup for all services controlled by this manager
	parent *cgroup2.Manager
}

// LoadAll Walks the root service path and enumerates all services and there respective groups
func (manager *Manager) LoadAll() error {
	qualifiedPath := ServicePath + "/" + manager.namespace

	_, err := os.Stat(qualifiedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			slog.Info("Service path does not exist, creating this", "path", qualifiedPath)

			if err := os.MkdirAll(qualifiedPath, 0600); err != nil {
				return fmt.Errorf("failed to create service path: %v", err)
			}

			return nil // if the path didn't exist, then there are no services to enumerate
		}

		return fmt.Errorf("unknown error when trying to create service path: %v", err)
	}

	if err := filepath.Walk(qualifiedPath, func(path string, info fs.FileInfo, err error) error {
		if path != qualifiedPath {
			if !info.IsDir() && strings.HasSuffix(path, ".json") {
				service, err := FromFile(path)
				if err != nil { // returning error here will stop reading any other found files
					slog.Error("failed to read service file, will be skipped", "err", err, "path", path)
					return nil
				}

				superv, err := NewSupervisor(service, manager.parent)
				if err != nil {
					slog.Error("failed to create supervisor for service, will be skipped", "err", err, "path", path)
					return nil
				}

				manager.services = append(manager.services, superv)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	slog.Info("enumerated service files for namespace", "ns", manager.namespace, "count", len(manager.services))

	return nil
}

// StartAll Starts all services the manager has enumerated in priority order
func (manager *Manager) StartAll(ctx context.Context) {
	manager.engine = task.New(manager.services)

	if err := manager.engine.PrioritizedWait(func(elem *Supervisor) error {
		elem.Start(ctx)
		return nil
	}); err != nil {
		slog.Error("Got error when trying to start all services", "err", err) // will never hit
	}
}

// StopAll Stops all services the manager has enumerated in priority order
//
// TODO: This should be reverse priority order
func (manager *Manager) StopAll(ctx context.Context) {
	if err := manager.engine.PrioritizedWait(func(elem *Supervisor) error {
		elem.Stop(ctx)
		return nil
	}); err != nil {
		slog.Error("Got error when trying to stop all services", "err", err)
	}
}

// New Initializes a new service.Manager and returns a pointer to it
func New(namespace string) (*Manager, error) {
	parent, err := cgroup2.NewManager(CgroupMountPoint, fmt.Sprintf("/%s.slice", namespace), new(cgroup2.Resources))
	if err != nil {
		return nil, err
	}

	return &Manager{
		namespace: namespace,
		services:  make([]*Supervisor, 0),
		parent:    parent,
	}, nil
}
