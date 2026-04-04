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

	// rootPhase The first phase that is executed after phase.MountPhase is executed
	rootPhase phase.Phase
}

func (init *Init) Start() {
	err := phase.NewMountPhase().Execute()
	if err != nil {
		log.Println("Root phase failed to execute!", err)
	}

	init.config = config.Default()
	init.config.AutomaticBootParams()

	init.log.Println("Enumerated the following boot arguments: ", init.config.BootParams())
	if init.config.Debug == true {
		init.log.Println("Got 'debug' in boot arguments. Starting debug phase...")

		err := phase.NewDebugPhase().Execute()
		if err != nil {
			init.log.Println("Debug phase failed to execute with: ", err)
		}
	}
}

func New(rootPhase phase.Phase) *Init {
	return &Init{
		log:       log.New(os.Stdout, "initd: ", log.Lshortfile|log.LstdFlags),
		rootPhase: rootPhase,
	}
}
