package root

import (
	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/factory"
	"github.com/spf13/cobra"
)

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd(factory *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servicectl",
		Short: "A simple service manager written in Go",
		Long:  "A simple service manager written in Go",
	}

	return cmd
}
