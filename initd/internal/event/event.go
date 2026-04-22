package event

import (
	"context"
	"log"

	"github.com/stevezaluk/immutable-os/initd/internal/acpi"
	"github.com/stevezaluk/immutable-os/initd/internal/reaper"
)

// Loop Core event loop for init. Blocks on powerEvent until an ACPI event is received. Once an ACPI event is received,
// the process reaper is started and all processes are killed before the root filesystem is unmounted. Finally, syscall.Reboot
// is called to pass control from the init system back over to the kernel
func Loop(ctx context.Context) {
	acpiSig := make(chan bool)

	go acpi.Wait(ctx, acpiSig)

	for {
		select {
		case <-acpiSig:
			log.Println("ACPI event received. Initiating shutdown")
			err := reaper.ReapProcesses()
			if err != nil {
				log.Fatalf("Error reaping processes: %v", err)
			}
			return
		case <-ctx.Done():
			return
		}
	}
}
