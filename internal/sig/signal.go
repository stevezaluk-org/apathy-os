package sig

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// SignalRegistration Represents a signal that should be caught and reacted to
type SignalRegistration struct {
	// C Channel that signals will be received at
	C chan os.Signal

	// Sig The signal registered
	Sig syscall.Signal

	// handlers Any handlers that need to be executed when this signal is received
	handlers []SignalHandlerFunc
}

// Execute Notify SignalRegistration.C when a signal is received and execute its list of handlers
func (registration *SignalRegistration) Execute(ctx context.Context) {
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

// NewRegistration Creates a new SignalRegistration and returns a pointer to it
func NewRegistration(sig syscall.Signal, handlers []SignalHandlerFunc) *SignalRegistration {
	return &SignalRegistration{
		C:        make(chan os.Signal, 1),
		Sig:      sig,
		handlers: handlers,
	}
}
