package sig

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"syscall"
)

// HandlerFunc A function to execute when a signal is received
type HandlerFunc func()

// Handler Handles taking specific actions when a Linux signal is received
type Handler struct {
	// C channel that signal events are sent to
	C chan os.Signal

	// mu Mutex lock for Handler.handlers
	mu sync.Mutex

	// signals All registered signals
	signals []*Registration
}

// RegisterSignal Creates a new SignalRegistration and registers it with the handler. Error is returned if you pass
// no signal handler functions. Go-routine safe
func (handler *Handler) RegisterSignal(sig syscall.Signal, f ...HandlerFunc) {
	handler.mu.Lock()
	defer handler.mu.Unlock()

	if len(f) == 0 {
		f = make([]HandlerFunc, 0)
	}

	handler.signals = append(handler.signals, NewRegistration(sig, f))
}

// HandleSignals Spawns a new go-routine per-unique signal that has been registered with the signal handler. If a signal
// has multiple functions assigned to it, then they are executed sequentially in the order that they are received. This
// function blocks until each spawned go-routine is finished. The context passed here is passed to each go-routine spawned
// for signal handling, so canceling it ends all signal processing.
func (handler *Handler) HandleSignals(ctx context.Context) {
	if len(handler.signals) == 0 {
		return
	}

	wg := new(sync.WaitGroup)
	for _, handler := range handler.signals {
		slog.Info("starting signal handler", "nr", handler.Sig, "name", handler.Sig.String())
		wg.Go(func() {
			handler.Execute(ctx)
		})
	}

	slog.Info("waiting for signals")
	wg.Wait()
}

// New Creates a new signal handler and returns a pointer to it
func New() *Handler {
	return &Handler{
		C:       make(chan os.Signal, 1),
		signals: make([]*Registration, 0),
	}
}
