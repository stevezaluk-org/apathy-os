package bootarg

import (
	"os"
	"reflect"
	"strings"
)

const (
	BootArgumentsFile      = "/proc/cmdline"
	BootArgumentsStructTag = "bootarg"
)

// BootArguments An abstraction over boot arguments available under /proc/cmdline. Only some common boot arguments
// are exposed as direct struct members here, however other boot arguments can be enumerated using the Raw method
type BootArguments struct {
	// bootArgs Raw boot arguments parsed from /proc/cmdline
	bootArgs map[string]string

	// Quiet A common boot argument for suppressing dmesg. Exposed here as well to control userland logging
	Quiet bool `bootarg:"quiet"`

	// Splash A common boot argument for displaying the splash screen during boot
	Splash bool `bootarg:"splash"`

	// Root The `root` boot argument which usually defines the device containing the root filesystem
	Root string `bootarg:"root"`

	// RootType The `rootfstype` boot argument which usually defines the filesystem type of the root filesystem
	RootType string `bootarg:"rootfstype"`
}

// Raw Returns raw boot arguments parsed from /proc/cmdline
func (bootArguments *BootArguments) Raw() map[string]string {
	return bootArguments.bootArgs
}

// Read Reads boot arguments from /proc/cmdline and returns map[string]string representing them
func Read() (map[string]string, error) {
	ret := make(map[string]string)

	contents, err := os.ReadFile(BootArgumentsFile)
	if err != nil {
		return nil, err
	}

	for _, bootarg := range strings.Split(string(contents), " ") {
		if strings.Contains(bootarg, "=") {
			split := strings.SplitN(bootarg, "=", 2)

			key := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])

			ret[key] = value

			continue
		}

		ret[bootarg] = ""
	}

	return ret, nil
}

// Unmarshal Parses raw boot arguments into an instance of BootArguments. `=` joined arguments (i.e. root=blah) are
// unmarshalled to strings, while lone boot arguments (i.e. quiet) are unmarshalled to bool's. Use the `bootarg` struct
// tag to control unmarshalling
func Unmarshal(raw map[string]string, v any) error {
	s := reflect.TypeOf(v)

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)

		if field.Tag.Get("boot") != "" {
			param, ok := raw[field.Tag.Get(BootArgumentsStructTag)]
			if ok {
				// We check to see if this is empty here as boot parameters are sometimes passed with no values. Ex:
				// params: ro quiet splash
				// This effectively lets us represent them as boolean's in Config
				if param == "" { // we got a bool
					if field.Type == reflect.TypeOf(true) {
						reflect.ValueOf(v).Elem().Field(i).SetBool(true)
					}

					continue
				}

				if field.Type == reflect.TypeOf(param) { // map all other values to their respective types
					reflect.ValueOf(v).Elem().Field(i).Set(reflect.ValueOf(param))
				}
			}
		}
	}

	return nil
}

// New Reads boot arguments from /proc/cmdline and returns an instantiated pointer to a BootArguments struct. Alias for
// calling Read followed by Unmarshal
func New() (*BootArguments, error) {
	raw, err := Read()
	if err != nil {
		return nil, err
	}

	ret := new(BootArguments)
	if err := Unmarshal(raw, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
