package mount

import (
	"errors"

	"github.com/stevezaluk/immutable-os/go-manifest/pkg/engine"
)

var ErrMountFailed = errors.New("mount failed")

// Handler Controls the processing of a list of mount points
type Handler struct {
	// engine The concurrency engine that the Handler uses
	engine *engine.Engine

	// mounts Mount list binned by priority
	mounts []*Mount
}

func (handler *Handler) MountAll(remount bool) error {
	return nil
}

// NewHandler Initialize a new mount.Handler and return a pointer to it
func NewHandler(mounts []*Mount) *Handler {
	return &Handler{mounts: mounts}
}
