package event

import (
	"log"

	"github.com/stevezaluk/container-os/init/internal/acpi"
)

// Loop Core event loop for init. Blocks on powerEvent until an ACPI event is received. Once an ACPI event is received,
// the process reaper is started and all processes are killed before the root filesystem is unmounted. Finally, syscall.Reboot
// is called to pass control from the init system back over to the kernel
func Loop() {
	powerEvent := make(chan bool)

	go acpi.Wait(powerEvent)
	for {
		select {
		case <-powerEvent:
			log.Println("ACPI event received")
		}
	}
}
