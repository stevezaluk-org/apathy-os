package main

import (
	"os"

	"github.com/stevezaluk-org/apathy-os/internal/cmd/serviced/factory"
	"github.com/stevezaluk-org/apathy-os/internal/cmd/serviced/root"
)

func main() {
	if err := root.NewRootCmd(factory.New()).Execute(); err != nil {
		os.Exit(1)
	}
}
