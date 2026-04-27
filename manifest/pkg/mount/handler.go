package mount

import (
	"errors"
	"log"
	"maps"
	"slices"

	internalErrors "github.com/stevezaluk/immutable-os/manifest/internal/errors"
	"github.com/stevezaluk/immutable-os/manifest/pkg/config"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sys/unix"
)

var ErrMountFailed = errors.New("mount failed")

// Handler Controls the processing of a list of mount points
type Handler struct {
	// mounts Mount list binned by priority
	mounts map[int]config.MountConfig
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
	priorities := slices.Sorted(maps.Keys(handler.mounts)) // sort by priority

	if remount {
		slices.Reverse(priorities)
	}

	for _, priority := range priorities {
		mountPoints := handler.mounts[priority]

		wg := new(errgroup.Group)

		// TODO: Need to reverse mountPoints here for `remount`
		for _, mountPoint := range mountPoints {
			log.Println("mounting: ", mountPoint.MountPoint)
			flags := uintptr(0)
			if remount {
				if mountPoint.ReadOnly { // no point in remounting an r/o filesystem as r/o
					continue
				}

				flags |= unix.MS_REMOUNT
				mountPoint.ReadOnly = true
			}

			wg.Go(func() error {
				if err := handler.mount(mountPoint, flags); err != nil {
					return err
				}

				return nil
			})
		}

		// Wait for all mounts in this priority class to complete
		if err := wg.Wait(); err != nil {
			return err
		}

	}

	return nil
}

// NewHandler Initialize a new mount.Handler and return a pointer to it
func NewHandler(mounts config.MountConfig) *Handler {
	return &Handler{mounts: mounts.Bin()}
}
