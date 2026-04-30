package service

import (
	"os"

	"github.com/stevezaluk/immutable-os/go-manifest/pkg/engine"
	"golang.org/x/sys/unix"
)

type Handler struct {
	engine *engine.Engine[*ProcessService]

	// processServices Any basic process based services that should be started
	processServices []*ProcessService
}

func (handler *Handler) startProcess(process *ProcessService) error {
	args := make([]string, 0)
	if process.Exec {
		args = append(args, process.Binary)
		args = append(args, process.Args...)

		if err := unix.Exec(process.Binary, args, os.Environ()); err != nil {
			return err
		}
	}

	return nil
}

// StartAll Starts all required services
func (handler *Handler) StartAll() error {
	if err := handler.engine.PrioritizedWait(handler.startProcess); err != nil {
		return err
	}

	return nil
}

// NewHandler Initializes a new service.Handler and returns a pointer to it
func NewHandler(processServices []*ProcessService) *Handler {
	return &Handler{processServices: processServices, engine: engine.New(processServices)}
}
