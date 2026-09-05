package config

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config Universal configuration options for go-serviced
type Config struct {
	// Socket The unix path to expose the daemons socket at
	Socket string `json:"socket" yaml:"socket" mapstructure:"socket"`

	// ServicePath The default unix path that service files should be checked at
	ServicePath string `json:"service_path" yaml:"service_path" mapstructure:"service_path"`

	// CgroupFilesystem The mount-point of the cgroup V2 filesystem
	CgroupFilesystem string `json:"cgroup_fs" yaml:"cgroup_fs" mapstructure:"cgroup_fs"`
}

// BindFlags A wrapper around viper.BindPFlags that provides access to the viper instance that the config
// structure keeps track of
func (config *Config) BindFlags(cmd *cobra.Command) error {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}

// Load Reads a config file present at the path, and loads them into Viper. If an error occurs when the file is read
// then Load will fall back to environment variables prefixed with `SERVICED`. Additionally, the path provided is
// sanitized and will have occurrences of `~` replaced with whatever os.Getenv(HOME) evaluates to
func (config *Config) Load(path string) error {
	if strings.HasPrefix(path, "~") {
		path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	}

	viper.SetConfigFile(path) // path is the config file; type inferred from ext

	if err := viper.ReadInConfig(); err != nil {
		viper.SetEnvPrefix("SERVICED")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
		viper.AutomaticEnv()
	}

	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}

// Default Initialize a new Config with sane-defaults and return a pointer to it
func Default() *Config {
	return &Config{
		Socket:           "/var/run/go-serviced/sock.socket",
		ServicePath:      "/usr/share/service.d",
		CgroupFilesystem: "/sys/fs/cgroup",
	}
}
