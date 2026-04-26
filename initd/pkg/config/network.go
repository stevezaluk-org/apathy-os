package config

// NetworkInterfaceConfig Defines how you want to configure network interface cards (NIC)
type NetworkInterfaceConfig struct {
	// Name The name of the network interface to configure
	Name string `yaml:"name"`

	// DHCP If set to true, acquire an IP address for this NIC through DHCP
	DHCP bool `yaml:"dhcp"`

	// IPAddress A static IP address to assign to the NIC. Cannot coexist with NetworkInterfaceConfig.Dhcp == true
	IPAddress string `yaml:"ipAddress"`

	// Gateway Your networks gateway/router
	Gateway string `yaml:"gateway"`

	// VLANTag The VLAN tag that you want to assign to this network interface
	VLANTag int `yaml:"vlanTag"`
}

// NetworkConfig Configuration block for defining your networking settings
type NetworkConfig struct {
	// Interfaces A list of network interfaces to configure. If left empty, then attempt to automatically identify one
	Interfaces []*NetworkInterfaceConfig `yaml:"interfaces"`

	// DNSServers Any DNS Servers
	DNSServers []string `yaml:"dnsServers"`

	// Hostname The hostname to assign the machine
	Hostname string `yaml:"hostname"`
}
