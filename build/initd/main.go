package main

import cmd "github.com/stevezaluk-org/apathy-os/internal/cmd/initd"

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		panic(err)
	}
}
