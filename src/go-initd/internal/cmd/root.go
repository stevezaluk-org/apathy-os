/*
Copyright © 2026 Steven A. Zaluk <szaluk@stevenzaluk.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"log"
	"os"

	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/manifest"
	"github.com/immutable-os/immutable-os/src/go-manifest/pkg/runner"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
)

var cfgFile string

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "init",
		Short: "go-initd is an init system and declarative system configuration tool for Unix like systems",
		Long:  `go-initd is an init system and declarative system configuration tool for Unix like systems.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if os.Getpid() != 1 {
				log.Fatal("go-initd must be launched as PID 1")
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Version: ", Version())

			earlyConfig, err := manifest.Read(viper.GetString("manifest"))
			if err != nil {
				log.Fatalf("got fatal error triyng to read early boot config: %v", err)
			}

			if err := runner.New(earlyConfig).Go(); err != nil {
				log.Fatalf("got fatal error trying to process early boot manfiest: %v", err)
			}

			if err := unix.Exec(
				viper.GetString("switch-root"),
				[]string{viper.GetString("switch-root"), viper.GetString("rootfs"), "/init"},
				os.Environ(),
			); err != nil {
				log.Fatalf("got fatal error when trying to switch_root: %v", err)
			}
		},
	}

	rootCmd.AddCommand(NewVersionCmd())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "set the config file for go-initd")
	rootCmd.Flags().StringP("manifest", "m", "/usr/share/manifest.d/init.yaml", "The manifest file to run for system initialization")
	rootCmd.Flags().Bool("system", false, "If set to true, put go-initd in system mode. Used during late boot and assumes a mounted rootfs")

	/*
		The flags below are a result of the current state of go-initd execution and are subject to deprecation soon. init
		should not be dependent on externalized tools initializing the system.
	*/
	rootCmd.Flags().String("rootfs", "/newroot", "The directory to switch root to after system initialization")
	rootCmd.Flags().String("busybox", "/usr/sbin/busybox", "The location of busybox to use for filesystem discovery")
	rootCmd.Flags().String("switch-root", "/usr/sbin/switch_root", "The location of the switch_root binary to execute during switch_root")

	if err := viper.BindPFlags(rootCmd.Flags()); err != nil {
		log.Fatalf("got fatal error trying to bind viper flags: %v", err)
	}

	return rootCmd
}
