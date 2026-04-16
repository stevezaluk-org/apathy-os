package acpi

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"syscall"

	"github.com/stevezaluk/container-os/init/internal/constants"
)

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// Wait An event loop that persistently reads for ACPI events from event handlers exposed at /dev/input/event0. This
// function is pretty crude at the moment and repeatedly loops with a blocking read on the path above. The real problem
// here is that we are not dynamically identifying which devices map to which event handlers so there is a chance that
// on other systems, /dev/input/event0 will be mapped to a different device entirely, and we won't actually receive an ACPI
// event. This will be rectified in a future commit
func Wait() {
	for {
		fd, err := os.Open(constants.ACPIEventPath)
		if err != nil {
			log.Fatal("Got error trying to open ACPI event file: ", err)
		}

		buffer := make([]byte, 4096)

		n, err := fd.Read(buffer)
		if err != nil {
			log.Fatal("Got error trying to read ACPI event file: ", err)
		}

		event := &InputEvent{}

		err = binary.Read(bytes.NewBuffer(buffer[:n]), binary.LittleEndian, event)
		if err != nil {
			log.Fatal("Got error trying to marshal ACPI event: ", err)
		}

		log.Println("Event read: ", event)
	}
}
