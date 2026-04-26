package mount

import "github.com/stevezaluk/immutable-os/initd/pkg/config"

type Handler struct {
	// config Defines the desired configuration for mounting filesystems
	config *config.MountConfig
}
