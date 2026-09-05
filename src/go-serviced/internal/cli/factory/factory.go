package factory

import "github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"

var (
	cachedConfig = config.Default()
)

type Factory struct {
	Config func(string) (*config.Config, error)
}

func New() *Factory {
	return &Factory{
		Config: func(path string) (*config.Config, error) {
			if cachedConfig != nil {
				return cachedConfig, nil
			}

			if err := cachedConfig.Load(path); err != nil {
				return nil, err
			}

			return cachedConfig, nil
		},
	}
}
