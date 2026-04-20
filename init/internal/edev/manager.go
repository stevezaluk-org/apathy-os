package edev

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	EVIOCGNAME uintptr = 0x80ff4506
	EVIOCGID   uintptr = 0x80084502
)

// DeviceManager Abstraction for enumerating and interacting with event devices
type DeviceManager struct {
	Devices []*Device
}

// getName Enumerates the name of the event device and returns it as a string. If the system call this method makes fails
// it will be returned as a wrapped error
func (manager *DeviceManager) getName(fd uintptr) (string, error) {
	// the unix.SYS_IOCTL trap always expects 256 bits here so we need to initialize a buffer of exactly that size
	var buffer [256]byte

	// We need to make a raw system call here, since the 'unix' and 'syscall' packages do not give us a direct way
	// of retrieving the name of an event device. This is a little bit convoluted and is not the most performant as
	// we need to make a system call for each device that we are enumerating. In the future, this should be deprecated,
	// and we should just enumerate from /proc/bus/input/devices
	_, _, err := unix.Syscall(
		unix.SYS_IOCTL,
		fd,
		EVIOCGNAME,
		uintptr(unsafe.Pointer(&buffer[0])),
	)
	if err != 0 {
		return "", fmt.Errorf("failed to make syscall with exit code: %d", err)
	}

	// ensure that we properly parse the populated byte slice
	return string(buffer[:bytes.IndexByte(buffer[:], 0)]), nil
}

// Enumerate Iterates through each event handler under /dev/input and enumerates basic information about the device. This
// includes the name, bus, product, version, and vendor. This function call is bound to be re-worked in the future as
// it makes an excessive amount of raw system calls to enumerate each device. This can instead be implemented by reading
// /proc/bus/input/devices which would drastically reduce the number of system calls that are being made. This really
// isn't the end of the world at the moment, as this only needs to happen once on system start and can be run concurrently
// to not block other start up routines.
func (manager *DeviceManager) Enumerate() error {
	err := filepath.Walk("/dev/input", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), "event") {
			fmt.Println("/dev/input/" + info.Name())
			file, err := os.OpenFile("/dev/input/"+info.Name(), unix.O_RDONLY|unix.O_NONBLOCK, 0)
			if err != nil {
				return err
			}

			devName, err := manager.getName(file.Fd())
			if err != nil {
				return err
			}

			manager.Devices = append(manager.Devices, &Device{
				Name:    devName,
				Handler: "/dev/input/" + info.Name(),
			})
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// NewManager Initializes a new device manager and returns a pointer to it.
func NewManager() *DeviceManager {
	return &DeviceManager{
		Devices: make([]*Device, 0),
	}
}
