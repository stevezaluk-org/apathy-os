package listcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stevezaluk-org/apathy-os/internal/cmd/serviced/factory"
	v1 "github.com/stevezaluk-org/apathy-os/internal/service/proto/service/v1"
)

func NewListCmd(factory *factory.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists all registered services",
		Long:  "Lists all registered services",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := factory.Client()
			if err != nil {
				return err
			}

			resp, err := client.List(context.Background(), &v1.ListServiceRequest{})
			if err != nil {
				return err
			}

			fmt.Println("Services:")
			for _, svc := range resp.Service {
				fmt.Printf("%s (%s)", svc.Name, svc.Exec)
			}

			return nil
		},
	}

	return cmd
}
