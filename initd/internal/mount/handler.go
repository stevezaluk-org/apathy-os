package mount

import (
	"fmt"
	"slices"

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

// mount Mounts a mount point defined in Handler.config.Mounts
func (handler *Handler) mount(mountPoint *config.MountPoint, flags ...uintptr) error {
	source, err := handler.findDevice(mountPoint)
	if err != nil {
		return err
	}

	mountFlags := uintptr(0)
	if len(flags) != 0 {
		for _, flag := range flags {
			mountFlags |= flag
		}
	}

	if mountPoint.ReadOnly {
		mountFlags |= unix.MS_RDONLY
	}

	// always implicitly trust the user knows what there doing and that filesystem types are accurate
	if err := unix.Mount(source, mountPoint.MountPoint, mountPoint.FilesystemType, mountFlags, ""); err != nil {
		return fmt.Errorf("failed to mount filesystem %s: %v", source, err)
	}

	return nil
}

// MountAll Processes all user and system defined mounts. By default, system mounts are processed unless remountRo is set
// to true. In that case, user mounts are processed first.
func (handler *Handler) MountAll(remountRo bool) error {
	// Mount points are generally processed in system first order, unless we are in system shutdown (and we are remounting
	// as read only). If that is the case we process user mounts first
	mountPoints := slices.Concat(handler.config.System, handler.config.User)
	if remountRo {
		slices.Reverse(mountPoints)
		unix.Sync() // Required to ensure that we don't miss any uncommited writes
	}

	// TODO: May need to add a "failFast" field in config.Mounts to determine if additional mounts should be processed if one fails
	// TODO: I still think a priority field would make sense here to determine the order in which mounts are processed
	for _, value := range mountPoints {
		flags := uintptr(0)
		if remountRo {
			// if remountRo is set, that means we are in system shutdown, and we gain nothing by trying to remount an already
			// mounted read only filesystem as read only. Would be completely wasting system calls
			if value.ReadOnly {
				continue
			}

			flags = unix.MS_REMOUNT
		}

		if err := handler.mount(value, flags); err != nil {
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
