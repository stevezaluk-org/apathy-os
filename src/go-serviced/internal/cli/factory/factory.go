package factory

import "github.com/immutable-os/immutable-os/src/go-serviced/pkg/config"

type Factory struct {
	Config func(string) (*config.Config, error)
}

func New() *Factory {
	return &Factory{
		Config: func(path string) (*config.Config, error) {
			c := config.Default()
			if err := c.Load(path); err != nil {
				return nil, err
			}

			return c, nil
		},
	}
}
