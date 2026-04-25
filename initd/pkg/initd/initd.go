package initd

import (
	"log"

	"github.com/stevezaluk/immutable-os/initd/internal/config"
	"github.com/stevezaluk/immutable-os/initd/internal/edev"
)

// Init Abstraction representing PID 1
type Init struct {
	// config Boot arguments read from /proc/cmdline.
	config *config.Config

	// devices Event Device Manager used for receiving ACPI events
	devices *edev.DeviceManager

	// log Basic STDOUT logger
	log *log.Logger
}

// New Instantiates an Init struct and returns a pointer to it
func New() *Init {
	return &Init{}
}
