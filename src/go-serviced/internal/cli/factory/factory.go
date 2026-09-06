package factory

import (
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	cachedConfig = config.Default()
)

type Factory struct {
	Config func() (*config.Config, error)
	Client func() (v1.ServicesClient, error)
}

func New() *Factory {
	factory := &Factory{
		Config: func() (*config.Config, error) {
			if cachedConfig != nil {
				return cachedConfig, nil
			}

			if err := cachedConfig.Load(viper.GetString("config")); err != nil {
				return nil, err
			}

			return cachedConfig, nil
		},
	}

	factory.Client = func() (v1.ServicesClient, error) {
		c, err := factory.Config()
		if err != nil {
			return nil, err
		}

		client, err := grpc.NewClient("unix://"+c.Socket,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		return v1.NewServicesClient(client), nil
	}

	return factory
}
