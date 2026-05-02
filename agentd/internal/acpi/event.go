package acpi

import (
	"bytes"
	"encoding/binary"

	"github.com/mdlayher/genetlink"
)

// Event Describes an ACPI event received from a generic netlink socket. The layout of this struct needs to match
// exactly what the kernel sends in raw binary through the generic netlink socket, which is why we are using explicitly
// sized byte slices
type Event struct {
	// DeviceClass The device that the event was received from
	DeviceClass [20]byte

	// BusId The kernel's Bus ID for the device
	BusId [15]byte

	// _ Padding to satisfy binary.Read
	_ [1]byte

	// Type The type of event received
	Type uint32

	// Data The data of the event
	Data uint32
}

// UnmarshalEvent Unmarshal's a generic netlink message into an acpi.Event message and returns a pointer to it.
func UnmarshalEvent(msg genetlink.Message) (*Event, error) {
	ret := new(Event)

	// binary.NativeEndian ensures that we only ever read in the host systems endianess
	if err := binary.Read(bytes.NewBuffer(msg.Data), binary.NativeEndian, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
