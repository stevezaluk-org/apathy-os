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

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "switch_root",
		Short: "A go implementation of the `switch_root` command from util-linux",
		Long:  "A go implementation of the `switch_root` command from util-linux",
	}

	cmd.AddCommand(NewVersionCmd())

	cmd.Flags().String("root", "/newroot", "The directory to switch_root to")
	cmd.Flags().String("init", "/init", "An absolute path to init")

	return cmd
}
