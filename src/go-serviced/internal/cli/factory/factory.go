package factory

import (
	"github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"
	v1 "github.com/immutable-os/immutable-os/src/go-serviced/proto/service/v1"
	"github.com/spf13/viper"
)

var (
	cachedConfig = config.Default()
)

type Factory struct {
	Config func() (*config.Config, error)
	Client func() (*v1.ServicesClient, error)
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

	return factory
}
