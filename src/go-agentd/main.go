package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/stevezaluk/immutable-os/src/go-agentd/internal/agent"
	"github.com/stevezaluk/immutable-os/src/go-manifest/pkg/manifest"
)

// TODO: Better CLI here
func main() {
	if os.Getpid() != 1 {
		panic(errors.New("agentd must be launched as PID 1"))
	}

	config, err := manifest.Read("/usr/share/manifest.d/system.yaml")
	if err != nil {
		panic(err)
	}

	log.Println(config)

	agentd := agent.New(config)
	if err := agentd.Start(); err != nil {
		panic(fmt.Errorf("failed to start agentd: %w", err))
	}
}
