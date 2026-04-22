package acpi

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"

	"github.com/stevezaluk/immutable-os/init/internal/constants"
	"github.com/stevezaluk/immutable-os/init/internal/edev"
)

const (
	KeyPower uint16 = 116
	EvKey    uint16 = 1
)

// readEvent Marshals an input event and checks to make sure that a received event is actually an ACPI event
//
// Deprecated: Use the functions available under the 'edev' package
func readEvent() (bool, error) {
	fd, err := os.Open(constants.ACPIEventPath)
	if err != nil {
		log.Fatal("Got error trying to open ACPI event file: ", err)
	}

	buffer := make([]byte, 4096)

	n, err := fd.Read(buffer)
	if err != nil {
		log.Fatal("Got error trying to read ACPI event file: ", err)
	}

	event := new(edev.InputEvent)

	err = binary.Read(bytes.NewBuffer(buffer[:n]), binary.LittleEndian, event)
	if err != nil {
		return false, err
	}

	if event.Type != EvKey && event.Code != KeyPower {
		return false, fmt.Errorf("event was not an ACPI event (Type: %d, Code %d)", event.Type, event.Code)
	}

	defer fd.Close()

	return true, nil
}

// Wait An event loop that persistently reads for ACPI events from event handlers exposed at /dev/input/event0. This
// function is pretty crude at the moment and repeatedly loops with a blocking read on the path above. The real problem
// here is that we are not dynamically identifying which devices map to which event handlers so there is a chance that
// on other systems, /dev/input/event0 will be mapped to a different device entirely, and we won't actually receive an ACPI
// event. This will be rectified in a future commit
//
// Deprecated: Use the functions available under the 'edev' package
func Wait(ctx context.Context, sig chan<- bool) {
	for {
		select {
		default:
			check, err := readEvent()
			if err != nil {
				log.Println("Got error trying to read ACPI event file: ", err)
			}

			// if we get an ACPI event then notify the parent goroutine so that we can initialize shutdown
			if check {
				sig <- true
				return
			}
		case <-ctx.Done():
			return
		}
	}
}
