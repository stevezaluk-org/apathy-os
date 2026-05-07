package service

import "sync"

// Supervisor An abstraction for managing the lifecycle of services. Provides core centralized functionality for evaluating
// the health of services defined in /etc/service.d/*.json.
type Supervisor struct {
	// processList Universal process list for tracking running/stopped services
	processList map[string]*Service

	// mu Mutex lock protecting concurrent writes to the process list
	mu sync.RWMutex
}

// NewSupervisor Initializes a new Supervisor and returns a pointer to it
func NewSupervisor() *Supervisor {
	return &Supervisor{
		processList: make(map[string]*Service),
		mu:          sync.RWMutex{},
	}
}
