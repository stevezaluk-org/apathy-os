package mount

import (
	"errors"
	"fmt"
	"slices"

	internalErrors "github.com/stevezaluk/immutable-os/manifest/internal/errors"
	"github.com/stevezaluk/immutable-os/manifest/pkg/config"
	"golang.org/x/sys/unix"
)

var ErrMountFailed = errors.New("mount failed")

// Handler Controls the processing of a list of mount points
type Handler struct {
	// mountPoints A list of mount points to process
	mounts config.MountConfig
}

// findDevice Returns an absolute path to a given block device
func (handler *Handler) findDevice(mount config.MountPoint) (string, error) {
	// we are not validating if `mount`. Device is actually populate here as we are trusting that config.MountConfig
	// was validated before it was passed to us
	if mount.FilesystemUUID == "" {
		return mount.Device, nil
	}

	device, err := FindFilesystem(mount.FilesystemUUID)
	if err != nil {
		return "", err
	}

	return device, nil
}

// mount Processes a single mount and returns a wrapped error if the mount fails. Merges any flags passed in the `flags`
// argument with unix.MS_RDONLY if the given config.MountPoint is read only
func (handler *Handler) mount(mount config.MountPoint, flags ...uintptr) error {
	device, err := handler.findDevice(mount)
	if err != nil {
		return err
	}

	mountFlags := uintptr(0)
	if len(flags) != 0 {
		for _, flag := range flags {
			mountFlags |= flag
		}
	}

	if mount.ReadOnly {
		mountFlags |= unix.MS_RDONLY
	}

	if err := unix.Mount(device, mount.MountPoint, mount.FilesystemType, mountFlags, ""); err != nil {
		return internalErrors.Wrap(ErrMountFailed, err)
	}

	return nil
}

// MountAll Mounts all mount points passed to the handler. MountAll will sort the mount point slice provided based on
// the config.MountPoint.Priority field and will concurrently mount where applicable. If remount is set to true, then
// each mount point will be remounted as read only in reverse order according to its priority
func (handler *Handler) MountAll(remount bool) error {
	if remount {
		// reverse our pre-sorted list so that we can remount in reverse priority order
		slices.Reverse(handler.mounts)
		unix.Sync()
	}

	for _, mount := range handler.mounts {
		flags := uintptr(0)
		if remount {
			// if remount is set, and we are processing a read-only mount, then there is no point in attempting to remount
			// a readonly filesystem as readonly. Doing so would harm performance and would waste system calls.

			// This also depends on system state not deviating from what is defined. If a defined read only mount point
			// was remounted as r/w, then we are potentially incurring data loss here
			if mount.ReadOnly {
				continue
			}

			flags = unix.MS_REMOUNT
		}

		if err := handler.mount(mount, flags); err != nil {
			return err
		}
	}

	return nil
}

// NewHandler Initialize a new mount.Handler and return a pointer to it
func NewHandler(mounts config.MountConfig) *Handler {
	return &Handler{mounts: mounts}
}
