package root

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stevezaluk-org/apathy-os/internal/cmd/serviced/daemoncmd"
	"github.com/stevezaluk-org/apathy-os/internal/cmd/serviced/factory"
	"github.com/stevezaluk-org/apathy-os/internal/cmd/serviced/listcmd"
	"github.com/stevezaluk-org/apathy-os/internal/machinery/cmdutil/versioncmd"
)

// NewRootCmd Initializes a new root command and returns a pointer to it
func NewRootCmd(factory *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "servicectl",
		Short: "A simple service manager written in Go",
		Long:  "A simple service manager written in Go",
	}

	versioncmd.NewVersionCmd(cmd)
	cmd.AddCommand(daemoncmd.NewDaemonCmd(factory))
	cmd.AddCommand(listcmd.NewListCmd(factory))

	cmd.PersistentFlags().StringP("config", "c", os.Getenv("HOME")+"/.go-serviced.yaml", "Config file for go-serviced (default: $HOME/.go-serviced.yaml)")
	_ = viper.BindPFlags(cmd.PersistentFlags())

	return cmd
}
