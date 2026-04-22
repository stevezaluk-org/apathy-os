package edev

import (
	"bytes"
	"encoding/binary"
	"syscall"
)

// InputEvent Represents an input event that would be received from event device
type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}

// NewEvent Unmarshal data from buffer into an InputEvent
func NewEvent(buffer []byte) (*InputEvent, error) {
	ret := new(InputEvent)

	err := binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
