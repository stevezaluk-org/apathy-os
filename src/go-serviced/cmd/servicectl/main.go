package main

import (
	"os"

	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/factory"
	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/root"
)

func main() {
	if err := root.NewRootCmd(factory.New()).Execute(); err != nil {
		os.Exit(1)
	}
}
