package agent

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/immutable-os/immutable-os/src/go-initd/internal/acpi"

	"github.com/immutable-os/immutable-os/src/go-initd/internal/reaper"
	"github.com/immutable-os/immutable-os/src/go-initd/internal/service"
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/manifest"
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/runner"
	"golang.org/x/sys/unix"
)

// SystemInit An abstraction around the agent as a whole. Manages the state of late PID 1 and service supervisor
type SystemInit struct {
	// config go-manifest config file
	config *manifest.Manifest

	// runner go-manifest runner used for initializing the system during startup
	runner *runner.Runner

	// supervisor Service supervisor for Agent
	supervisor *service.Supervisor

	// sigChild Channel responsible for delivering SIGCHLD signals. Used for reaping processes
	sigChild chan os.Signal

	// sigTerm Channel responsible for delivering SIGTERM events, an alternative way of initiating a shutdown
	sigTerm chan os.Signal

	// sigAcpi Channel responsible for receiving ACPI events, indicating that the system should shutdown
	sigAcpi chan bool
}

// EventLoop Starts the core agent event loop
// TODO: Both the process reaper and service shutdown is occurring at the same time.
func (systemInit *SystemInit) EventLoop() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Establish notifications for the signals that we need to catch
	signal.Notify(systemInit.sigChild, unix.SIGCHLD)
	signal.Notify(systemInit.sigTerm, unix.SIGTERM, unix.SIGINT)

	// start our ACPI event loop
	go acpi.EventLoop(ctx, systemInit.sigAcpi)
	go systemInit.supervisor.EventLoop(ctx)

	for {
		select {
		case <-systemInit.sigChild:
			log.Println("got sigchld")
			// TODO: Handle wait for pid here
		case <-systemInit.sigAcpi:
			log.Println("got ACPI")
			cancel()
			if err := systemInit.Shutdown(); err != nil {
				return err
			}
		case <-systemInit.sigTerm:
			log.Println("got sigterm")
			cancel()
			if err := systemInit.Shutdown(); err != nil {
				return err
			}
		}
	}
}

// Shutdown Initiates shutdown logic
//
// TODO: Not remounting (reverse order) filesystems here. Leads to data loss
// TODO: No real service management here. Basic unix.Kill(-1)
func (systemInit *SystemInit) Shutdown() error {
	if err := reaper.ReapProcesses(); err != nil {
		return err
	}

	// actually shutdown the system
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		return err
	}

	return nil
}

// Start Initial entrypoint for the systemInit Processes system/user manifest files and then moves to the event loop
func (systemInit *SystemInit) Start() error {
	if err := systemInit.runner.Go(); err != nil {
		return err
	}

	supervisor, err := service.NewSupervisor(service.DefaultServicePath)
	if err != nil {
		return err
	}

	systemInit.supervisor = supervisor

	// blocks forever
	if err := systemInit.EventLoop(); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Agent and returns a pointer to it
func New(config *manifest.Manifest) *SystemInit {
	return &SystemInit{
		runner:   runner.New(config),
		sigChild: make(chan os.Signal, 1),
		sigTerm:  make(chan os.Signal, 1),
		sigAcpi:  make(chan bool),
	}
}
