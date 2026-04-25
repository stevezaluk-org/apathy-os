package mount

import (
	"fmt"
	"os"
	"syscall"

	"github.com/stevezaluk/immutable-os/initd/internal/constants"
)

// SwitchRoot Calls switch_root on /root and sets /sbin/init [late] as its new init process
func SwitchRoot(root string) error {
	err := syscall.Exec(
		constants.SwitchrootPath,
		[]string{constants.SwitchrootPath, root, constants.LateInitPath, "--late"},
		os.Environ(),
	)
	if err != nil {
		return fmt.Errorf("failed to switch_root to rootfs: %v", err)
	}

	return nil
}
