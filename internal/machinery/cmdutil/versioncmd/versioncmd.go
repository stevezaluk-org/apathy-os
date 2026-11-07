package versioncmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stevezaluk-org/apathy-os/internal/version"
)

// NewVersionCmd Initializes a new version command and attaches it as a sub-command to what you pass here
func NewVersionCmd(root *cobra.Command) {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Displays the semantic version, commit SHA, build time, and go runtime that this release was built against",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", version.String())
		},
	}

	root.AddCommand(versionCmd)
}
