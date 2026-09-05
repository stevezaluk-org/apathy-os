package root

import "github.com/spf13/cobra"

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servicectl",
		Short: "A simple service manager written in Go",
		Long:  "A simple service manager written in Go",
	}

	return cmd
}
