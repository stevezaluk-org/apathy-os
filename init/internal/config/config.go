package config

import (
	"os"
	"reflect"
	"strings"
)

const (
	BootParamsFile string = "/proc/cmdline"
)

// Config Global config structure for controlling init
type Config struct {
	// Debug If set to true, initramfs will drop to a temporary BusyBox shell before starting
	Debug bool `boot:"debug"`

	// Quiet Suppress all logging
	Quiet bool `boot:"quiet"`
}

// AutomaticBootParams Parse boot parameters from /proc/cmdline and merge them to Config
func (config *Config) AutomaticBootParams() {
	contents, err := os.ReadFile(BootParamsFile)
	if err != nil {
		panic(err) // failed to read boot params
	}

	bootParams := make(map[string]string)
	for _, param := range strings.Split(string(contents), " ") {
		if strings.Contains(param, "=") {
			kv := strings.SplitN(param, "=", 2)
			bootParams[kv[0]] = kv[1]

			continue
		}

		bootParams[param] = ""
	}

	s := reflect.TypeOf(*config)

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)

		if field.Tag.Get("boot") != "" {
			param, ok := bootParams[field.Tag.Get("boot")]
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
}

// Default Returns defaults for Config
func Default() *Config {
	return &Config{
		Debug: false,
		Quiet: false,
	}
}
