package daemon

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/containerd/cgroups/v2/cgroup2"
	"github.com/immutable-os/immutable-os/src/go-serviced/internal/log"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/api"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/supervisor"
	"golang.org/x/sys/unix"
)

type Daemon struct {
	// config The config file passed to the daemon at construction
	config *config.Config

	// services A list of services enumerated from config.ServicePath
	services map[string]*supervisor.Supervisor

	// cgroup A cgroup representing the system.slice
	cgroup *cgroup2.Manager

	// api The API server that this daemon is managing
	api *api.Api

	// signals Channel for receiving unix signals being sent to the daemon
	signals chan os.Signal
}

// Load Loads relevant service files from the path passed as an argument
//
// TODO: It would be nice if we could have an abstraction for a service "group". Would make this nice and thin
func (daemon *Daemon) Load() error {
	cgroup, err := cgroup2.NewManager(daemon.config.CgroupFilesystem, "/system.slice", &cgroup2.Resources{})
	if err != nil {
		return err
	}

	err = filepath.Walk(daemon.config.ServicePath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".json") {
			superv, err := supervisor.FromFile(path, cgroup)
			if err != nil {
				return err
			}

			slog.Info("Discovered service", "svc", superv.Service.Name)
			daemon.services[superv.Service.Name] = superv
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// handleSignals Listens for unix signals sent to the daemon and executes actions when they are received
func (daemon *Daemon) handleSignals(ctx context.Context) {
	signal.Notify(daemon.signals, unix.SIGHUP, unix.SIGTERM, unix.SIGINT)

	for {
		switch {
		case <-daemon.signals == unix.SIGINT: // this needs to stop services
			if err := daemon.api.CleanUp(); err != nil {
				slog.Error("failed to clean up resources (SIGINT): ", "err", daemon.api.CleanUp())
			}
		case <-daemon.signals == unix.SIGTERM: // this needs to stop services
			if err := daemon.api.CleanUp(); err != nil {
				slog.Error("failed to clean up resources (SIGTERM): ", "err", daemon.api.CleanUp())
			}
		case <-daemon.signals == unix.SIGHUP: // this needs to reload service files
		}
	}
}

// Start Starts the daemon service and
func (daemon *Daemon) Start(ctx context.Context) error {
	sigCtx, sigCancel := context.WithCancel(ctx)
	defer sigCancel()

	go daemon.handleSignals(sigCtx)

	// blocks the current go-routine
	if err := daemon.api.Serve(); err != nil {
		return err
	}

	return nil
}

// NewDaemon Initializes a new daemon and returns a pointer to it. Also creates the `system.slice` cgroup
func NewDaemon(config *config.Config) (*Daemon, error) {
	daemon := &Daemon{
		config:   config,
		services: make(map[string]*supervisor.Supervisor),
		signals:  make(chan os.Signal, 1),
	}

	if err := daemon.Load(); err != nil {
		return nil, err
	}

	grpcServer, err := api.New(config, &daemon.services)
	if err != nil {
		return nil, err
	}
	daemon.api = grpcServer
	defer log.Init()

	return daemon, nil
}
