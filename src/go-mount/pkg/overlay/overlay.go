package overlay

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

const (
	// tmpPrefix Directory prefix used when Overlay.UpperDir or Overlay.WorkDir == ""
	tmpPrefix = "/tmp/overlay"
)

type Overlay struct {
	// UpperDir The upper directory of the overlayfs
	UpperDir string `yaml:"upperdir" json:"upperdir"`

	// LowerDir The lower directory of the overlayfs. Usually the mount point itself
	LowerDir string `yaml:"lowerdir" json:"lowerdir"`

	// WorkDir A writable working directory for staged writes to be performed in.
	WorkDir string `yaml:"workdir" json:"workdir"`

	// Merged The merged directory of the overlayfs
	Merged string `yaml:"merged" json:"merged"`
}

// Mount Mounts the overlay filesystem. If upperdir and workdir are omitted, then subdirectories under /tmp are used
// instead. This effectively makes the mount read-only with changes being wiped when /tmp is unmounted
func Mount(overlay *Overlay) error {
	if overlay.UpperDir == "" {
		upper := "/tmp/overlay/upperdir/" + overlay.UpperDir

		if err := os.MkdirAll(upper, 0600); err != nil {
			return err
		}

		overlay.UpperDir = upper
	}

	if overlay.WorkDir == "" {
		work := "/tmp/overlay/workdir/" + overlay.WorkDir
		if err := os.MkdirAll(work, 0600); err != nil {
			return err
		}

		overlay.WorkDir = work
	}

	opts := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s",
		overlay.LowerDir, overlay.UpperDir, overlay.WorkDir)

	// source and fstype here are always overlay since we are laying this ontop of an existing mount
	if err := unix.Mount("overlay", overlay.Merged, "overlay", 0, opts); err != nil {
		return fmt.Errorf("failed to mount overlayfs: %v", err)
	}

	return nil
}
