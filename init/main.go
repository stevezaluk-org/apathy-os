/*
Copyright © 2026 Steven A. Zaluk
*/
package main

import (
	"log"

	"github.com/stevezaluk/container-os/init/internal/config"
	"github.com/stevezaluk/container-os/init/internal/debug"
	"github.com/stevezaluk/container-os/init/internal/mount"
)

func main() {
	log.SetPrefix("initd ")

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
}
