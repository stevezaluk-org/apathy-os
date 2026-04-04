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

	// rootPhase The first phase that is executed before Init.config is populated. Usually phase.MountPhase
	rootPhase phase.Phase
}

func (init *Init) Start() {
	err := init.rootPhase.Execute()
	if err != nil {
		log.Println("Root phase failed to execute!", err)
	}

	init.config = config.Default()
	init.config.AutomaticBootParams()

	init.log.Println("Boot arguments: ", init.config.BootParams())
}

func New(rootPhase phase.Phase) *Init {
	return &Init{
		log:       log.New(os.Stdout, "initd: ", log.Lshortfile|log.LstdFlags),
		rootPhase: rootPhase,
	}
}
