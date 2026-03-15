/*
Copyright © 2026 Steven A. Zaluk
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "init",
	Short: "Init system and service manager for container-os",
	Long: `init acts as the PID 1 for container-os and is responsible for starting core services
and mounting the root filesystem. Similar to systemd, init can be started in system mode or user
mode. System mode is executed when the initramfs is mounted (before the root file system is mounted)
and init is started in user mode after the root filesystem has been mounted.

When init starts, boot arguments are enumerated to determine the mode that it should initially start in.
If the boot argument CONTAINER_OS_INIT=y is passed, then init is started in system mode and will automatically
transition into user mode. Additionally, the flags --system and --user and provided as override flags.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("system", "s", false, "Starts init in system mode")
	rootCmd.Flags().BoolP("user", "s", false, "Starts init in user mode")
}
