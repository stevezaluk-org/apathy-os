package main

import "github.com/stevezaluk-org/apathy-os/internal/cmd/initd/rootcmd"

func main() {
	if err := rootcmd.NewRootCmd().Execute(); err != nil {
		panic(err)
	}
}
