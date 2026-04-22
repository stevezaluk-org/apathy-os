/*
Copyright © 2026 Steven A. Zaluk
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"

	"github.com/stevezaluk/immutable-os/init/internal/config"
	"github.com/stevezaluk/immutable-os/init/internal/debug"
	"github.com/stevezaluk/immutable-os/init/internal/event"
	"github.com/stevezaluk/immutable-os/init/internal/mount"
	"github.com/stevezaluk/immutable-os/init/internal/shell"
)

func usage() {
	fmt.Println("Usage: initd [late] [early] ARGS")
	fmt.Println("\tearly - Handles the early boot process of mounting virtual filesystems, establishing an early debug shell, and mounting the rootfs")
	fmt.Println("\tlate - Handles the late boot process of establishing networking, starting other processes, and processing shutdown events")

	os.Exit(1)
}

func early() {
	log.SetPrefix("initd/early ")

	log.Printf("Init has started. Starting virtual filesystem mount")
	err := mount.MountVirtualFS()
	if err != nil {
		log.Fatal("got fatal error when trying to mount virtual filesystem: ", err)
	}

	// now that we have /proc mounted we can parse boot arguments
	globalConfig, err := config.New()
	if err != nil {
		log.Fatal("got fatal error when trying to read boot arguments: ", err)
	}

	// the symlinks created by this function call are required for mounting rootfs (blkid, findfs, and switch_root)
	if err := shell.CreateSymlinks(); err != nil {
		log.Fatal("got fatal error when trying to create symlinks: ", err)
	}

	err = mount.MountRootfs(globalConfig)
	if err != nil {
		log.Fatal("got fatal error when trying to mount rootfs: ", err)
	}

	err = mount.MoveVirtualMounts()
	if err != nil {
		log.Fatalf("moving virtual filesystems failed with: %v", err)
	}

	err = mount.SwitchRoot()
	if err != nil {
		log.Fatalf("failed to switch to rootfs: %v", err)
	}
}

func late() {
	log.SetPrefix("initd/late ")
	log.Printf("welcome to late init")

	globalConfig, err := config.New()
	if err != nil {
		log.Fatal("got fatal error when trying to read boot arguments: ", err)
	}

	// force into BusyBox shell if we get debug
	if globalConfig.Debug {
		log.Printf("'debug' was present in boot arguments. BusyBox shell will start shortly")
		err := debug.StartDebug()
		if err != nil {
			log.Fatal("got fatal error when trying to start debug: ", err)
		}
	}

	wg := new(sync.WaitGroup)
	wg.Go(func() { event.Loop(context.Background()) })
	wg.Wait()

	err = mount.RemountRootfs()
	if err != nil {
		log.Fatalf("got fatal error when trying to remove rootfs: %v", err)
	}

	err = syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		log.Fatalf("got fatal error when trying to reboot: %v. You should force a shutdown", err)
	}
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 2 {
		usage()
	}

	switch os.Args[1] {
	case "early":
		early()
	case "late":
		late()
	default:
		usage()
	}

}
