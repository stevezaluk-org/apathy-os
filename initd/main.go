/*
Copyright © 2026 Steven A. Zaluk
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stevezaluk/immutable-os/initd/pkg/initd"
)

func usage() {
	fmt.Println("Usage: initd [flags]")
	fmt.Println("Args:")
	fmt.Println("\t-h, --help - Shows this help message")
	fmt.Println("Phases:")
	fmt.Println("\t--early - Handles the early boot process of mounting virtual filesystems, and mounting the rootfs")
	fmt.Println("\t--late - Handles the late boot process of establishing networking, starting other processes, and processing shutdown events")
	fmt.Println("\t--event-loop - Start initd's event loop and listen for ACPI events and signals")
	fmt.Println("\t--shutdown - Reap processes, remount rootfs as read only, and shutdown")

	os.Exit(1)
}

func main() {
	flagSet := flag.NewFlagSet("initd", flag.ExitOnError)
	flagSet.Usage = usage

	early := flagSet.Bool("early", false, "")
	late := flagSet.Bool("late", false, "")
	eventLoop := flagSet.Bool("event-loop", false, "")
	shutdown := flagSet.Bool("shutdown", false, "")

	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	init := initd.New()
	if *early {
		if err := init.Early(); err != nil {
			log.Fatalf("Got fatal error from early: %v", err)
		}
	}

	if *late {
		if err := init.Late(); err != nil {
			log.Fatalf("Got fatal error from late: %v", err)
		}
	}

	if *eventLoop {
		if err := init.EventLoop(); err != nil {
			log.Fatalf("Got fatal error from event-loop: %v", err)
		}
	}

	if *shutdown {
		if err := init.Shutdown(); err != nil {
			log.Fatalf("Got fatal error from shutdown: %v", err)
		}
	}
}
