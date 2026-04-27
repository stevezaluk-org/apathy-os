package bootarg

// BootArguments An abstraction over boot arguments available under /proc/cmdline
type BootArguments struct {
	// bootArgs Raw boot arguments parsed from /proc/cmdline
	bootArgs map[string]string
}

// Raw Returns raw boot arguments parsed from /proc/cmdline
func (bootArguments *BootArguments) Raw() map[string]string {
	return bootArguments.bootArgs
}

// New Initialize a new BootArguments structure and return a pointer to it
func New() *BootArguments {
	return &BootArguments{}
}
