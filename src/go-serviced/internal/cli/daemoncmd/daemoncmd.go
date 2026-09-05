package daemoncmd

import (
	"context"
	"fmt"

	"github.com/immutable-os/immutable-os/src/go-serviced/internal/cli/factory"
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/daemon"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewDaemonCmd Initializes a new daemon command and returns a pointer to it
func NewDaemonCmd(factory *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "daemon",
		Short: "Starts the go-serviced daemon",
		Long: `The go-serviced daemon is what servicectl interacts with when starting/stopping/loading new services. When
the daemon is started, the file path defined under config.Config.ServicePath is read for service files, and are loaded
into the daemons internal service list and are started. 

The v1.Service.Priority field indicates the order that services will be started in. Lower numbers means the service will 
start sooner, higher numbers will start services later, and services with the same priority will start concurrently.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := factory.Config(viper.GetString("config"))
			if err != nil {
				return fmt.Errorf("could not load config: %v", err)
			}

			if err := c.BindFlags(cmd); err != nil {
				return fmt.Errorf("could not bind flags: %v", err)
			}

			d, err := daemon.NewDaemon(c)
			if err != nil {
				return fmt.Errorf("could not initialize daemon: %v", err)
			}

			if err := d.Start(context.Background()); err != nil {
				return fmt.Errorf("could not start daemon: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().String("socket", "/var/run/go-serviced/sock.socket", "File path to the socket to expose gRPC at")
	cmd.Flags().String("service_path", "/usr/share/service.d", "File path to look for sevrices under")
	cmd.Flags().String("cgroup_fs", "/sys/fs/cgroup", "Mount path to a filesystem of type `cgroup2`. Usually `/sys/fs/cgroup`")

	return cmd
}
