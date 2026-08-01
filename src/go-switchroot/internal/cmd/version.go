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
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	// SemanticVersion The semantic version of go-switchroot. Should match the tag that this was built against
	SemanticVersion string = "undefined"

	// Commit The short git commit SHA that go-switchroot was built against
	Commit string = "undefined"

	// BuildTime The date/time that go-switchroot was built
	BuildTime string = "undefined"
)

// Version Returns the go-switchroot version string
func Version() string {
	return fmt.Sprintf("go-switchroot %s (commit %s, built at %s, runtime %s)",
		SemanticVersion, Commit, BuildTime, runtime.Version())
}

// NewVersionCmd Initializes a new version command and returns a pointer to it
func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version string for go-switchroot",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version())
		},
	}
}
