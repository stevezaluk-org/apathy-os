package mount

import (
	"fmt"

	"github.com/stevezaluk/immutable-os/initd/pkg/config"
	"golang.org/x/sys/unix"
)

// Handler Processes the config.Mounts section. Holds the responsibility of processing user/system mounts
type Handler struct {
	// config Defines the desired configuration for mounting filesystems
	config *config.MountConfig
}

// findDevice Finds the actual device path under /dev for a mount point
func (handler *Handler) findDevice(mountPoint *config.MountPoint) (string, error) {
	if mountPoint.FilesystemUUID == "" { // always value filesystem UUIDs over raw device paths as they are more stable
		return mountPoint.Device, nil
	}

	device, err := FindFilesystem(mountPoint.FilesystemUUID)
	if err != nil {
		return "", fmt.Errorf("failed to find filesystem: %v", err)
	}

	return device, nil
}

// Mount Mounts a mount point defined in Handler.config.Mounts
func (handler *Handler) Mount(mountPoint *config.MountPoint) error {
	source, err := handler.findDevice(mountPoint)
	if err != nil {
		return err
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

// Go Process all system and user mounts
// TODO: May need to add a "failFast" field in config.Mounts to determine if additional mounts should be processed if one fails
func (handler *Handler) Go() error {
	// system mounts are always handled first in the event that user mounts depend on them
	// TODO: Adding a priority field here may be useful to control the order of how mounts are processed
	for _, value := range handler.config.System {
		if err := handler.Mount(value); err != nil {
			return err
		}
	}

	for _, value := range handler.config.User {
		if err := handler.Mount(value); err != nil {
			return err
		}
	}

	return nil
}

// NewHandler Initializes a mount.Handler and returns a pointer to it
func NewHandler(config *config.MountConfig) *Handler {
	return &Handler{
		config: config,
	}
}
