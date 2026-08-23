package user

import "fmt"

// LinuxUser Represents a Linux User that should be created during system boot
type LinuxUser struct {
	// Username The username
	Username string `yaml:"username"`

	// HomeDirectory The users home directory
	HomeDirectory string `yaml:"homeDir"`

	// Password A password hash that gets written to /etc/shadow
	Password string `yaml:"password"`

	// UID The ID of the user itself
	UID int `yaml:"uid"`

	// GID The ID of the users primary group
	GID int `yaml:"gid"`

	// Comment The GECOS of the user. Can be an abritrary string
	Comment string `yaml:"comment"`

	// Shell The default shell for the user
	Shell string `yaml:"shell"`
}

// String Returns the user metadata formatted as a string so that it can be written to /etc/passwd
func (linuxUser *LinuxUser) String() string {
	return fmt.Sprintf("%s:x:%d:%d:%d:%s:%s",
		linuxUser.Username, linuxUser.UID, linuxUser.GID, linuxUser.HomeDirectory, linuxUser.Comment, linuxUser.Shell)
}

// NewUser Initializes a new user with the set username and returns a pointer to it
func NewUser(username string, shell string) *LinuxUser {
	return &LinuxUser{Username: username, Shell: shell}
}
