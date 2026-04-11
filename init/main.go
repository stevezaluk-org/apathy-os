/*
Copyright © 2026 Steven A. Zaluk
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stevezaluk/container-os/init/internal/config"
	"github.com/stevezaluk/container-os/init/internal/debug"
	"github.com/stevezaluk/container-os/init/internal/mount"
	"github.com/stevezaluk/container-os/init/internal/shell"
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

	// force into BusyBox shell if we get debug
	if globalConfig.Debug {
		log.Printf("'debug' was present in boot arguments. BusyBox shell will start shortly")
		err := debug.StartDebug()
		if err != nil {
			log.Fatal("got fatal error when trying to start debug: ", err)
		}
	}

	err = mount.MountRootfs(globalConfig)
	if err != nil {
		log.Fatal("got fatal error when trying to mount rootfs: ", err)
	}
}

func late() {
	log.SetPrefix("initd/late ")
	log.Printf("welcome to late init")

	rootls, err := shell.RunCmd("/usr/bin/ls", "/")
	if err != nil {
		log.Fatal("got fatal error when trying to list rootfs: ", err)
	}

	log.Printf("rootfs contents: %v", rootls)
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
