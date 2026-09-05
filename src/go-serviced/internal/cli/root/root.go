package root

import (
	"os"

	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/daemoncmd"
	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/factory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd(factory *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servicectl",
		Short: "A simple service manager written in Go",
		Long:  "A simple service manager written in Go",
	}

	cmd.AddCommand(daemoncmd.NewDaemonCmd(factory))

	cmd.PersistentFlags().StringP("config", "c", os.Getenv("HOME")+"/.go-serviced.yaml", "Config file for go-serviced (default: $HOME/.go-serviced.yaml)")
	_ = viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}
