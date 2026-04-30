package edev

import (
	"context"
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

	n, err := fd.Read(data)
	if err != nil {
		return nil, err
	}

	event, err := NewEvent(data[:n])
	if err != nil {
		return nil, err
	}
	return event, nil
}

// Watch Streams received input events to Device.C. Performs a basic blocking read and should eventually be migrated
// over to inotify. If an error is received, then the go-routine is canceled and the function returns
func (device *Device) Watch(ctx context.Context) {
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
}

// NewDevice Constructs a new device and returns a handler to it
func NewDevice(name string, handler string) *Device {
	return &Device{
		Name:    name,
		Handler: handler,
		C:       make(chan *InputEvent),
	}
}
