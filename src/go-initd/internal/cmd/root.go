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
	"github.com/spf13/cobra"
)

var cfgFile string

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "init",
		Short: "go-initd is an init system and declarative system configuration tool for Unix like systems",
		Long:  `go-initd is an init system and declarative system configuration tool for Unix like systems.`,
	}

	rootCmd.AddCommand(NewVersionCmd())

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "set the config file for go-initd")

	return rootCmd
}
