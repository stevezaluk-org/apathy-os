package mount

import "golang.org/x/sys/unix"

// Flags A map of recognized mount flags
var Flags = map[string]uintptr{
	"ro":       unix.MS_RDONLY,
	"rw":       unix.MS_RDONLY,
	"remount":  unix.MS_REMOUNT,
	"noatime":  unix.MS_NOATIME,
	"relatime": unix.MS_RELATIME,
	"noexec":   unix.MS_NOEXEC,
	"nodev":    unix.MS_NODEV,
	"nosuid":   unix.MS_NOSUID,
}
