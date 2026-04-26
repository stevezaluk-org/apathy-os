package mount

import (
	"fmt"

	"github.com/stevezaluk/immutable-os/initd/pkg/config"
	"golang.org/x/sys/unix"
)

type Handler struct {
	// config Defines the desired configuration for mounting filesystems
	config *config.MountConfig
}

// Mount Mounts a mount point defined in Handler.config.Mounts
func (handler *Handler) Mount(mountPoint *config.MountPoint) error {
	source := mountPoint.Device
	if mountPoint.FilesystemUUID != "" { // always value filesystem UUIDs over raw device paths as they are more stable
		device, err := FindFilesystem(mountPoint.FilesystemUUID)
		if err != nil {
			return fmt.Errorf("failed to find filesystem: %v", err)
		}

		source = device
	}

	// uintptr(0) == no mount flags
	mountFlags := uintptr(0)
	if mountPoint.ReadOnly {
		// additional mount flags can be passed to unix.Mount with the | operator
		mountFlags |= unix.MS_RDONLY
	}

	// always implicitly trust the user knows what there doing and that filesystem types are accurate
	if err := unix.Mount(source, mountPoint.MountPoint, mountPoint.FilesystemType, mountFlags, ""); err != nil {
		return fmt.Errorf("failed to mount filesystem %s: %v", source, err)
	}

	return nil
}

// NewHandler Initializes a mount.Handler and returns a pointer to it
func NewHandler(config *config.MountConfig) *Handler {
	return &Handler{
		config: config,
	}
}
