package initd

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/stevezaluk-org/apathy-os/internal/initd/acpi"
	"github.com/stevezaluk-org/apathy-os/internal/initd/reaper"
	"github.com/stevezaluk-org/apathy-os/internal/service/config"
	"github.com/stevezaluk-org/apathy-os/internal/service/daemon"
	"github.com/stevezaluk-org/apathy-os/internal/sig"
	"golang.org/x/sys/unix"
)

// SystemInit An abstraction around the agent as a whole. Manages the state of late PID 1 and service supervisor
type SystemInit struct {
	// reaper The process reaper for PID 1
	reaper *reaper.Reaper

	// signalHandler Signal handler for reacting to signals received by initd
	signalHandler *sig.Handler

	// serviceDaemon Daemon for controlling service execution
	serviceDaemon *daemon.Daemon

	// sigAcpi Channel responsible for receiving ACPI events, indicating that the system should shutdown
	sigAcpi chan bool
}

// EventLoop Starts the core agent event loop
// TODO: Both the process reaper and service shutdown is occurring at the same time.
// TODO: Bad context handling here
func (systemInit *SystemInit) EventLoop() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register the signals that PID1 is responsible for handling
	systemInit.signalHandler.RegisterSignal(unix.SIGCHLD, func() {
		if err := systemInit.reaper.Reap(); err != nil {
			slog.Error("Got error trying to reap proccess. This can indicate that zombie's are accumulating", "err", err)
			cancel()
		}
	})

	systemInit.signalHandler.RegisterSignal(unix.SIGTERM, func() {
		if err := systemInit.Shutdown(); err != nil {
			slog.Error("Got error when trying to shutdown the system", "err", err)
			cancel()
		}
	})

	// start our ACPI event loop
	go acpi.EventLoop(ctx, systemInit.sigAcpi)
	go systemInit.signalHandler.HandleSignals(ctx)
	go systemInit.serviceDaemon.Serve(ctx)

	for {
		select {
		case <-ctx.Done():
			cancel()
			return nil
		case <-systemInit.sigAcpi:
			log.Println("got ACPI power off event. Shutting down system")
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := systemInit.serviceDaemon.StopAll(ctx); err != nil {
		return err
	}

	if err := systemInit.reaper.Fallback(); err != nil {
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
	serviceDaemon, err := daemon.NewDaemon(config.Default())
	if err != nil {
		return err
	}

	systemInit.serviceDaemon = serviceDaemon

	// blocks forever
	if err := systemInit.EventLoop(); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Agent and returns a pointer to it
func New() *SystemInit {
	return &SystemInit{
		reaper:        reaper.New(),
		signalHandler: sig.New(),
		sigAcpi:       make(chan bool),
	}
}
