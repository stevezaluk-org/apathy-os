package initd

import (
	"log"
	"os"

	"github.com/stevezaluk/container-os/init/internal/config"
	"github.com/stevezaluk/container-os/init/internal/phase"
)

// Init An abstraction representing the initd system as a whole. Controls the entire initd execution flow, primarily
// through starting and stopping phases. Also holds the "Root" phase i.e. the phase that executes before anything else
// does, usually phase.MountPhase
type Init struct {
	log *log.Logger

	// config Global configuration values populated with /proc/cmdline after mounting occurs
	config *config.Config

	// executor Responsible for executing init phases
	executor *phase.Executor
}

// StartPhases Starts init phase execution
func (init *Init) StartPhases() {
	if err := init.executor.Start(); err != nil {
		init.log.Fatalf("Phase executor returned with the following error: %s", err.Error())
	}
}

func New(rootPhase phase.Phase) *Init {
	return &Init{
		log:      log.New(os.Stdout, "initd: ", log.Lshortfile|log.LstdFlags),
		executor: phase.NewExecutor(),
	}
}
