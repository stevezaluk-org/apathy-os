package mount

import "github.com/stevezaluk/immutable-os/manifest/pkg/config"

// Handler Controls the processing of a list of mount points
type Handler struct {
	// mountPoints A list of mount points to process
	mountPoints []*config.MountPoint
}

// NewHandler Initialize a new mount.Handler and return a pointer to it
func NewHandler() *Handler {
	return &Handler{mountPoints: make([]*config.MountPoint, 0)}
}
