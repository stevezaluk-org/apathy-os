package sig

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// Registration Represents a signal that should be caught and reacted to
type Registration struct {
	// C Channel that signals will be received at
	C chan os.Signal

	// Sig The signal registered
	Sig syscall.Signal

	// handlers Any handlers that need to be executed when this signal is received
	handlers []HandlerFunc
}

// Execute Notify Registration.C when a signal is received and execute its list of handlers
func (registration *Registration) Execute(ctx context.Context) {
	signal.Notify(registration.C, registration.Sig)

	for {
		select {
		case <-ctx.Done():
			return
		case <-registration.C:
			for _, f := range registration.handlers {
				f()
			}
		}
	}
}

// NewRegistration Creates a new Registration and returns a pointer to it
func NewRegistration(sig syscall.Signal, handlers []HandlerFunc) *Registration {
	return &Registration{
		C:        make(chan os.Signal, 1),
		Sig:      sig,
		handlers: handlers,
	}
}
