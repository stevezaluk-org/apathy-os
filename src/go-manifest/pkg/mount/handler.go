package mount

import (
	"errors"

	"github.com/stevezaluk/immutable-os/src/go-manifest/pkg/engine"
	"golang.org/x/sys/unix"
)

var ErrMountFailed = errors.New("mount failed")

// Handler Controls the processing of a list of mount points
type Handler struct {
	// engine The concurrency engine that the Handler uses
	engine *engine.Engine[*Mount]

	// mounts Mount list binned by priority
	mounts []*Mount
}

// remountFilesystem Remounts a filesystem as read only. Existing read only filesystems are not remounted
func (handler *Handler) remountFilesystem(mountPoint *Mount) error {
	if mountPoint.ReadOnly {
		// Return early here as to not attempt a remount on a read-only filesystem
		return nil
	}

	if err := handler.mountFilesystem(mountPoint, unix.MS_REMOUNT); err != nil {
		return err
	}

	return nil
}

// mountFilesystem Processes a Mount and mounts the filesystem
func (handler *Handler) mountFilesystem(mountPoint *Mount, flags ...uintptr) error {
	device, err := FindFilesystem(mountPoint)
	if err != nil {
		return err
	}

	// this can technically break if a caller passes unix.MS_RDONLY as a flag instead of setting it in the mount point
	mountFlags := uintptr(0)
	if len(flags) != 0 {
		for _, flag := range flags {
			mountFlags |= flag
		}
	}

	if mountPoint.ReadOnly {
		mountFlags |= unix.MS_RDONLY
	}

	// this function assumes that data is already validated before it makes its way to this system call
	if err := unix.Mount(device, mountPoint.MountPoint, mountPoint.FilesystemType, mountFlags, ""); err != nil {
		return err
	}

	return nil
}

// MountAll Entry point for the mount handler. If remount is set, then unix.MS_REMOUNT and unix.MS_RDONLY are passed
// to unix.Mount. Usually used as a part of system shutdown. Existing read only filesystems will not be remounted
func (handler *Handler) MountAll(remount bool) error {
	f := func(mountPoint *Mount) error {
		return handler.mountFilesystem(mountPoint)
	}
	if remount {
		f = func(mountPoint *Mount) error {
			return handler.remountFilesystem(mountPoint)
		} // this does not remount in reverse order
	}

	if err := handler.engine.PrioritizedWait(f); err != nil {
		return err
	}

	return nil
}

// NewHandler Initialize a new mount.Handler and return a pointer to it
func NewHandler(mounts []*Mount) *Handler {
	return &Handler{mounts: mounts, engine: engine.New(mounts)}
}
