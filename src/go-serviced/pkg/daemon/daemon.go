package daemon

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/containerd/cgroups/v2/cgroup2"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/supervisor"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

type Daemon struct {
	// config The config file passed to the daemon at construction
	config *config.Config

	// services A list of services enumerated from config.ServicePath
	services map[string]*supervisor.Supervisor

	// cgroup A cgroup representing the system.slice
	cgroup *cgroup2.Manager
}

// Load Loads relevant service files from the path passed as an argument
//
// TODO: It would be nice if we could have an abstraction for a service "group". Would make this nice and thin
func (daemon *Daemon) Load() error {
	cgroup, err := cgroup2.NewManager(daemon.config.CgroupFilesystem, "system.slice", nil)
	if err != nil {
		return err
	}

	err = filepath.Walk(daemon.config.ServicePath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			file, err := os.OpenFile(path, os.O_RDONLY, 0644)
			if err != nil {
				return err
			}

			data := make([]byte, info.Size())
			n, err := file.Read(data)
			if err != nil {
				return err
			}

			service := &v1.Service{}
			err = protojson.Unmarshal(data[:n], service)
			if err != nil {
				return err
			}

			superv, err := supervisor.New(service, cgroup)
			if err != nil {
				return err
			}

			log.Printf("Service queued for start: %s", path)
			daemon.services[superv.Service.Name] = superv
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// NewDaemon Initializes a new daemon and returns a pointer to it. Also creates the `system.slice` cgroup
func NewDaemon(config config.Config) (*Daemon, error) {
	daemon := &Daemon{
		config:   &config,
		services: make(map[string]*supervisor.Supervisor),
	}

	if err := daemon.Load(); err != nil {
		return nil, err
	}

	return daemon, nil
}
