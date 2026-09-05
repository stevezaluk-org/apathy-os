package main

import (
	"os"

	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/root"
)

func main() {
	if err := root.NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
