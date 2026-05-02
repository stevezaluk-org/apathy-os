package agent

import (
	"log"
	"os"
	"os/signal"

	"github.com/stevezaluk/immutable-os/agentd/internal/edev"
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

	// sigChild Channel responsible for delivering SIGCHLD signals. Used for reaping processes
	sigChild chan os.Signal

	// sigTerm Channel responsible for delivering SIGTERM events, an alternative way of initiating a shutdown
	sigTerm chan os.Signal

	// sigAcpi Channel responsible for receiving ACPI events, indicating that the system should shutdowwn
	sigAcpi chan *edev.InputEvent
}

// EventLoop Starts the core agent event loop
func (agent *Agent) EventLoop() error {

	// Establish notifications for the signals that we need to catch
	signal.Notify(agent.sigChild, unix.SIGCHLD)
	signal.Notify(agent.sigTerm, unix.SIGTERM, unix.SIGINT)

	for {
		select {
		case <-agent.sigChild:
			log.Println("got sigchdl")
		case <-agent.sigAcpi:
		case <-agent.sigTerm:
			log.Println("got sigterm/acpi event")
		}
	}
}

// Start Initial entrypoint for the agent. Processes system/user manifest files and then moves to the event loop
func (agent *Agent) Start() error {
	if err := agent.runner.Go(); err != nil {
		return err
	}

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
		sigAcpi:  make(chan *edev.InputEvent),
	}
}
