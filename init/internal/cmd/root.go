/*
Copyright © 2026 Steven A. Zaluk
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/stevezaluk/container-os/init/internal/config"
	"github.com/stevezaluk/container-os/init/internal/initd"
)

var globalConfig *config.Config

var rootCmd = &cobra.Command{
	Use:   "initd",
	Short: "Init system and service manager for container-os",
	Long: `initd acts as the PID 1 for container-os and is responsible for starting core services
and mounting the root filesystem. Similar to systemd, initd can be started in system mode or user
mode. System mode is executed when the initramfs is mounted (before the root file system is mounted)
and initd is started in user mode after the root filesystem has been mounted.

When initd starts, boot arguments are enumerated to determine the mode that it should initially start in.
If the boot argument CONTAINER_OS_INIT=y is passed, then initd is started in system mode and will automatically
transition into user mode. Additionally, the subcommands 'system' and 'user' can be called to invoke this behavior
directly`,
	Run: func(cmd *cobra.Command, args []string) {
		i := initd.New(nil)
		i.StartPhases()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
