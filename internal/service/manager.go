package service

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

const (
	ServicePath = "/usr/share/service.d"
)

// Manager Global manager for all services running on the systems. Each manager is associated with one or more services, named by its
// namespace. Each manager gets a cgroup allocated to it, with each service under the managers control getting its own child cgroup.
type Manager struct {
	// namespace The namespace for this manager
	namespace string

	// services A slice of services that the manager has enumerated
	services []*Service
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

				manager.services = append(manager.services, service)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	slog.Info("enumerated service files for namespace", "ns", manager.namespace, "count", len(manager.services))

	return nil
}

// New Initializes a new service.Manager and returns a pointer to it
func New(namespace string) *Manager {
	return &Manager{
		namespace: namespace,
	}
}
