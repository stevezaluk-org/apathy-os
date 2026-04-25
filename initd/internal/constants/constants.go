package constants

const (
	// BlkidPath The path to blkid on initramfs
	BlkidPath string = "/sbin/blkid"

	// FindfsPath The path to findfs on initramfs
	FindfsPath string = "/sbin/findfs"

	// SwitchrootPath The path to switch_root on initramfs
	SwitchrootPath string = "/sbin/switch_root"

	// BusyBoxPath The path to busybox on initramfs
	BusyBoxPath string = "/usr/sbin/busybox"
)

const (
	// BootParamsFile Unix standard file path on where boot parameters are exposed
	BootParamsFile string = "/proc/cmdline"

	// TemporaryRootMount Absolute path for where rootfs is temporarily mounted before switch_root is invoked
	TemporaryRootMount string = "/root"

	// LateInitPath Absolute path to where init is located on the rootfs
	LateInitPath string = "/usr/sbin/init"
)
