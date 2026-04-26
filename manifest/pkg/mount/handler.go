package mount

import (
	"errors"

	"github.com/stevezaluk/immutable-os/manifest/pkg/config"
)

var ErrMountFailed = errors.New("mount failed")

// Handler Controls the processing of a list of mount points
type Handler struct {
	// mountPoints A list of mount points to process
	mounts config.MountConfig
}

// NewHandler Initialize a new mount.Handler and return a pointer to it
func NewHandler(mounts config.MountConfig) *Handler {
	return &Handler{mounts: mounts}
}
