package sig

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"syscall"
)

// SignalHandlerFunc A function to execute when a signal is received
type SignalHandlerFunc func()

// SignalHandler Handles taking specific actions when a Linux signal is received
type SignalHandler struct {
	// C channel that signal events are sent to
	C chan os.Signal

	// mu Mutex lock for SignalHandler.handlers
	mu sync.Mutex

	signals []*SignalRegistration
}

// RegisterSignal Creates a new SignalRegistration and registers it with the handler. Error is returned if you pass
// no signal handler functions. Go-routine safe
func (signalHandler *SignalHandler) RegisterSignal(sig syscall.Signal, f ...SignalHandlerFunc) error {
	if len(f) == 0 {
		return fmt.Errorf("must register at least one handler a signal")
	}

	signalHandler.mu.Lock()
	defer signalHandler.mu.Unlock()

	signalHandler.signals = append(signalHandler.signals, NewRegistration(sig, f))

	return nil
}

// HandleSignals Spawns a new go-routine per-unique signal that has been registered with the signal handler. If a signal
// has multiple functions assigned to it, then they are executed sequentially in the order that they are received. This
// function blocks until each spawned go-routine is finished. The context passed here is passed to each go-routine spawned
// for signal handling, so canceling it ends all signal processing.
func (signalHandler *SignalHandler) HandleSignals(ctx context.Context) {
	if len(signalHandler.signals) == 0 {
		return
	}

	wg := new(sync.WaitGroup)
	for _, handler := range signalHandler.signals {
		slog.Info("starting handler", "sig", handler.Sig, "f", handler.handlers)
		wg.Go(func() {
			handler.Execute(ctx)
		})
	}

	slog.Info("waiting for signals")
	wg.Wait()
}

// New Creates a new signal handler and returns a pointer to it
func New() *SignalHandler {
	return &SignalHandler{
		C:       make(chan os.Signal, 1),
		signals: make([]*SignalRegistration, 0),
	}
}
