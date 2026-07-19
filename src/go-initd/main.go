package main

import "github.com/immutable-os/immutable-os/src/go-initd/internal/cmd"

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		panic(err)
	}
}
