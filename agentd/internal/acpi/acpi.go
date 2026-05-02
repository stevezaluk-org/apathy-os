package acpi

import (
	"context"
	"log"
	"strings"

	"github.com/mdlayher/genetlink"
)

const (
	acpiFamily = "acpi_event"
)

// processEvents Returns true on the first power button event that is received
func processEvents(messages []genetlink.Message) (bool, error) {

	for _, msg := range messages {
		event, err := UnmarshalEvent(msg)
		if err != nil {
			continue // skip on an invalid message
		}

		if strings.Contains(string(event.DeviceClass[:]), "button/power") {
			return true, nil // we got a power button press, return true
		}
	}

	return false, nil
}

// EventLoop Loop forever until we receive a "button/power" event from a generic netlink socket. This indicates that a
// power button event was received, and we should start gracefully shutdown. Notify the channel passed in the second
// argument of this function.
func EventLoop(context context.Context, sig chan bool) {
	sock, err := genetlink.Dial(nil)
	if err != nil {
		log.Fatal(err)
	}

	fam, err := sock.GetFamily(acpiFamily)
	if err != nil {
		log.Fatal(err)
	}

	for _, group := range fam.Groups {
		if err := sock.JoinGroup(group.ID); err != nil {
			log.Fatal(err)
		}
	}

	for {
		select {
		case <-context.Done():
			return
		default:
			msg, _, err := sock.Receive()
			if err != nil {
				log.Fatal(err)
			}

			received, err := processEvents(msg)
			if err != nil {
				log.Fatal(err)
			}

			if received {
				sig <- true
				return // return on the first power event we receive so that we can initiate shutdown
			}
		}
	}
}
