package edev

import "syscall"

// InputEvent Represents an input event that would be received from event device
type InputEvent struct {
	Time  syscall.Timeval
	Type  uint16
	Code  uint16
	Value int32
}
