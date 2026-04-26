package mount

import "github.com/stevezaluk/immutable-os/initd/pkg/config"

type Handler struct {
	// config Defines the desired configuration for mounting filesystems
	config *config.MountConfig
}

// NewHandler Initializes a mount.Handler and returns a pointer to it
func NewHandler(config *config.MountConfig) *Handler {
	return &Handler{
		config: config,
	}
}
