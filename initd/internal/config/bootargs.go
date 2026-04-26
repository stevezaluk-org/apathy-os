package config

import (
	"os"
	"reflect"
	"strings"

	"github.com/stevezaluk/immutable-os/initd/internal/constants"
)

// BootArgs All boot arguments enumerated from /proc/cmdline
type BootArgs struct {
	// bootParams The raw boot arguments parsed from /proc/cmdline
	bootParams map[string]string

	// Debug If set to true, initramfs will drop to a temporary BusyBox shell before starting
	Debug bool `boot:"debug"`

	// Quiet Suppress all logging
	Quiet bool `boot:"quiet"`

	// Root A hint towards the root filesystem. Usually contains a string like UUID=<UUIDHERE> that represents the UUID of the filesystem
	Root string `boot:"root"`
}

// BootParams Returns raw boot arguments parsed from /proc/cmdline
func (bootArgs *BootArgs) BootParams() map[string]string {
	return bootArgs.bootParams
}

// ParseBootArgs Parse boot arguments from /proc/cmdline and merge them
func ParseBootArgs() (*BootArgs, error) {
	config := new(BootArgs)

	contents, err := os.ReadFile(constants.BootParamsFile)
	if err != nil {
		return nil, err
	}

	if config.bootParams == nil {
		config.bootParams = make(map[string]string)
	}

	for _, param := range strings.Split(string(contents), " ") {
		if strings.Contains(param, "=") {
			kv := strings.SplitN(param, "=", 2)
			config.bootParams[kv[0]] = kv[1]

			continue
		}

		config.bootParams[param] = ""
	}

	s := reflect.TypeOf(*config)

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)

		if field.Tag.Get("boot") != "" {
			param, ok := config.bootParams[field.Tag.Get("boot")]
			if ok {
				// We check to see if this is empty here as boot parameters are sometimes passed with no values. Ex:
				// params: ro quiet splash
				// This effectively lets us represent them as boolean's in Config
				if param == "" { // we got a bool
					if field.Type == reflect.TypeOf(true) {
						reflect.ValueOf(config).Elem().Field(i).SetBool(true)
					}

					continue
				}

				if field.Type == reflect.TypeOf(param) { // map all other values to there respective types
					reflect.ValueOf(config).Elem().Field(i).Set(reflect.ValueOf(param))
				}
			}
		}
	}

	return config, nil
}
