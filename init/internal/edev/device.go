package edev

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"os"
)

// Device Represents an event device
type Device struct {
	Name    string
	Handler string

	C chan *InputEvent
}

// blockingRead Performs a blocking read on the event handler and returns an InputEvent when one
// is received
func (device *Device) blockingRead() (*InputEvent, error) {
	fd, err := os.OpenFile(device.Handler, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	data := make([]byte, 4096)
	event := new(InputEvent)

	n, err := fd.Read(data)
	if err != nil {
		return nil, err
	}

	err = binary.Read(bytes.NewBuffer(data[:n]), binary.LittleEndian, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// Watch Start a go-routine that streams received input events to Device.C. As of the moment, this performs a basic
// blocking read and should eventually be migrated over to inotify
func (device *Device) Watch(ctx context.Context) error {
	if device.Handler == "" {
		return errors.New("device handler is empty")
	}
	go func() {
		for {
			select {
			default:
				event, err := device.blockingRead()
				if err != nil {
					return
				}

				device.C <- event
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
