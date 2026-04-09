package phase

import (
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

// NewExecutor Instantiates a new phase.Executor with VMount as its root phase. Also initializes a new logger under the
// prefix of phase/executor
func NewExecutor() *Executor {
	return &Executor{
		Root: NewMountPhase(),
		log:  internalLog.New("phase/executor"),
	}
}
