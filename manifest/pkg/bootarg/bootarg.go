package bootarg

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

// New Initialize a new BootArguments structure and return a pointer to it
func New() *BootArguments {
	return &BootArguments{}
}
