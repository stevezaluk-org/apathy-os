package main

import (
	"log"
	"os"

	"github.com/stevezaluk/immutable-os/go-manifest/pkg/manifest"
	"github.com/stevezaluk/immutable-os/go-manifest/pkg/runner"
)

// All constant paths are under initramfs, not rootfs
const (
	BusyBoxPath       = "/usr/sbin/busybox"
	EarlyBootManifest = "/usr/share/manifest.d/init.yaml"
)

// symlinks Initialize required symlinks for mounting root filesystem
func symlinks() error {
	required := []string{"findfs", "switch_root"}

	for _, link := range required {
		if err := os.Symlink(BusyBoxPath, "/usr/sbin/"+link); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if os.Getpid() != 1 {
		log.Fatal("initd must be launched as PID 1")
	}

	earlyConfig, err := manifest.Read(EarlyBootManifest)
	if err != nil {
		log.Fatalf("got fatal error triyng to read early boot config: %v", err)
	}

	if err := symlinks(); err != nil {
		log.Fatalf("got fatal error trying to establish symlinks for rootfs mount: %v", err)
	}

	if err := runner.New(earlyConfig).Go(); err != nil {
		log.Fatalf("got fatal error trying to process early boot manfiest: %v", err)
	}
}
