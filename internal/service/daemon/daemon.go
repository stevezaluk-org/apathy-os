package daemon

import (
	"context"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/containerd/cgroups/v2/cgroup2"
	"github.com/stevezaluk-org/apathy-os/internal/machinery/log"

	"github.com/stevezaluk-org/apathy-os/internal/api/serviced/api"
	"github.com/stevezaluk-org/apathy-os/internal/machinery/task"
	"github.com/stevezaluk-org/apathy-os/internal/service/config"
	"github.com/stevezaluk-org/apathy-os/internal/service/supervisor"
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

	// engine The engine responsible for concurrently starting/stopping services
	engine *task.Engine[*supervisor.Supervisor]
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
		case <-daemon.signals == unix.SIGINT || <-daemon.signals == unix.SIGTERM: // this needs to stop services
			if err := daemon.api.CleanUp(); err != nil {
				slog.Error("failed to clean up resources (SIGINT/SIGTERM): ", "err", daemon.api.CleanUp())
			}
		case <-daemon.signals == unix.SIGHUP: // this needs to reload service files
		}
	}
}

// StopAll Stops all services discovered by the daemon
func (daemon *Daemon) StopAll(ctx context.Context) error {
	err := daemon.engine.PrioritizedWait(func(superv *supervisor.Supervisor) error {
		slog.Info("Stopping service", "svc", superv.Service.Name)

		// TODO: This is kind of clunky and its exactly what we are doing at supervisor.Start
		cleanUpCtx, cleanUpCancel := context.WithTimeout(ctx, 20*time.Second)
		defer cleanUpCancel()

		if err := superv.Stop(cleanUpCtx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// StartAll Starts all services discovered by the daemon
func (daemon *Daemon) StartAll(ctx context.Context) error {
	err := daemon.engine.PrioritizedWait(func(superv *supervisor.Supervisor) error {
		svcCtx, cancel := context.WithCancelCause(ctx)

		slog.Info("Starting service", "svc", superv.Service.Name)
		if err := superv.Start(svcCtx, cancel); err != nil {
			return err
		}

		superv.Watch(svcCtx, cancel)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// Serve Starts the daemon service and any services that it has discovered
func (daemon *Daemon) Serve(ctx context.Context) error {
	sigCtx, sigCancel := context.WithCancel(ctx)
	defer sigCancel()

	go daemon.handleSignals(sigCtx)
	go daemon.StartAll(context.Background())

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
	daemon.engine = task.New[*supervisor.Supervisor](slices.Collect(maps.Values(daemon.services)))

	grpcServer, err := api.New(config, &daemon.services)
	if err != nil {
		return nil, err
	}
	daemon.api = grpcServer
	defer log.Init()

	return daemon, nil
}
