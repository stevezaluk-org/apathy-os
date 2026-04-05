/*
Copyright © 2026 Steven A. Zaluk
*/

package cmd

import (
	"github.com/spf13/cobra"
)

// systemCmd represents the system command
var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Start init in system mode",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)
}
