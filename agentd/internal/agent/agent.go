package agent

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/stevezaluk/immutable-os/agentd/internal/acpi"
	"github.com/stevezaluk/immutable-os/agentd/internal/reaper"
	"github.com/stevezaluk/immutable-os/agentd/internal/service"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/manifest"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/runner"
	"golang.org/x/sys/unix"
)

// Agent An abstraction around the agent as a whole. Manages the state of late PID 1 and service supervisor
type Agent struct {
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

	// sigAcpi Channel responsible for receiving ACPI events, indicating that the system should shutdowwn
	sigAcpi chan bool
}

// EventLoop Starts the core agent event loop
func (agent *Agent) EventLoop() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Establish notifications for the signals that we need to catch
	signal.Notify(agent.sigChild, unix.SIGCHLD)
	signal.Notify(agent.sigTerm, unix.SIGTERM, unix.SIGINT)

	// start our ACPI event loop
	go acpi.EventLoop(ctx, agent.sigAcpi)

	for {
		select {
		case <-agent.sigChild:
			log.Println("got sigchld")
			// TODO: Handle wait for pid here
		case <-agent.sigAcpi:
			log.Println("got ACPI")
			if err := agent.Shutdown(); err != nil {
				return err
			}
		case <-agent.sigTerm:
			log.Println("got sigterm")
			if err := agent.Shutdown(); err != nil {
				return err
			}
		}
	}
}

// Shutdown Initiates shutdown logic
//
// TODO: Not remounting (reverse order) filesystems here. Leads to data loss
// TODO: No real service management here. Basic unix.Kill(-1)
func (agent *Agent) Shutdown() error {
	if err := reaper.ReapProcesses(); err != nil {
		return err
	}

	// actually shutdown the system
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_POWER_OFF); err != nil {
		return err
	}

	return nil
}

// Start Initial entrypoint for the agent. Processes system/user manifest files and then moves to the event loop
func (agent *Agent) Start() error {
	if err := agent.runner.Go(); err != nil {
		return err
	}

	supervisor, err := service.NewSupervisor(service.DefaultServicePath)
	if err != nil {
		return err
	}

	agent.supervisor = supervisor

	// blocks forever
	if err := agent.EventLoop(); err != nil {
		return err
	}

	return nil
}

// New Initializes a new Agent and returns a pointer to it
func New(config *manifest.Manifest) *Agent {
	return &Agent{
		runner:   runner.New(config),
		sigChild: make(chan os.Signal, 1),
		sigTerm:  make(chan os.Signal, 1),
		sigAcpi:  make(chan bool),
	}
}
