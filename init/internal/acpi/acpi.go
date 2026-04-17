package acpi

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/stevezaluk/container-os/init/internal/constants"
)

const (
	KeyPower uint16 = 116
	EvKey    uint16 = 1
)

type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// checkEvent Marshals an input event and checks to make sure that a received event is actually an ACPI event
func checkEvent(buf []byte) (bool, error) {
	event := new(InputEvent)

	err := binary.Read(bytes.NewBuffer(buf), binary.LittleEndian, event)
	if err != nil {
		return false, err
	}

	if event.Type != EvKey && event.Code != KeyPower {
		return false, fmt.Errorf("event was not an ACPI event (Type: %d, Code %d)", event.Type, event.Code)
	}

	return true, nil
}

// Wait An event loop that persistently reads for ACPI events from event handlers exposed at /dev/input/event0. This
// function is pretty crude at the moment and repeatedly loops with a blocking read on the path above. The real problem
// here is that we are not dynamically identifying which devices map to which event handlers so there is a chance that
// on other systems, /dev/input/event0 will be mapped to a different device entirely, and we won't actually receive an ACPI
// event. This will be rectified in a future commit
func Wait(sig chan<- bool) {
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

		isAcpi, err := checkEvent(buffer[:n])
		if err != nil {
			log.Fatal("Got error trying to check ACPI event: ", err)
		}

		if isAcpi {
			log.Printf("ACPI event received. Notifying parent goroutine")
			sig <- true
		}
	}
}
