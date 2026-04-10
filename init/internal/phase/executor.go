package phase

import (
	"errors"
	"log"

	"github.com/stevezaluk/container-os/init/internal/config"
	internalLog "github.com/stevezaluk/container-os/init/internal/log"
)

// Executor A linked-list that is responsible for executing startup phases. A phase.Phase can represent a task that is
// required for starting the system
type Executor struct {
	// Root The first phase that gets executed. Always phase.vmount
	Root Phase

	// Config The config file marshaled from boot arguments found at /proc/cmdline
	Config *config.Config

	// log An instance of log.Logger used for propagating log events
	log *log.Logger
}

func (executor *Executor) initialize() error {
	// realistically we should never hit this, but we check anyway
	if executor.Root == nil {
		return errors.New("executor: root phase was nil")
	}

	// we always execute our first phase before we even iterate through the other phases as we want to ensure
	// that at least virtual filesystems are mounted before we continue

	// This is mandatory for a few reasons, most notably for enumerating configuration values from /proc/cmdline, but
	// also because there is really not much that can happen without virtual filesystems mounted
	if err := executor.Root.Execute(); err != nil {
		return err
	}
	// bump up execution to ensure that we don't execute the same phase twice
	executor.Root = executor.Root.Next()

	// after we have successfully mounted virtual filesystems we can call Config.AutomaticBootParams(). This will
	// automatically read /proc/cmdline and us reflection to hydrate executor.Config with configuration values.
	// Read more about this in config/config.go
	executor.Config = config.Default()
	executor.Config.AutomaticBootParams()

	if executor.Config.Debug {
		// catching debug in the boot arguments means the user wants to invoke the debug shell. Always executed
		// after virtual filesystem mounting

		// TODO: Maybe think about removing this as a phase
		if err := NewDebugPhase().Execute(); err != nil {
			return err
		}
	}

	return nil
}

// Start Run through phase execution and return any errors
func (executor *Executor) Start() error {
	executor.log.Printf("Starting phase execution pre-flight")
	if err := executor.initialize(); err != nil {
		return err
	}

	currentPhase := executor.Root
	for currentPhase != nil {
		if err := currentPhase.Execute(); err != nil {
			return err
		}

		currentPhase = currentPhase.Next()
	}

	return nil
}

// NewExecutor Instantiates a new phase.Executor with VMount as its root phase. Also initializes a new logger under the
// prefix of phase/executor
func NewExecutor() *Executor {
	return &Executor{
		Root: NewMountPhase(),
		log:  internalLog.New("phase/executor"),
	}
}
