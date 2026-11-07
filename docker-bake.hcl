// TAG The tag that gets assigned to a built image. Defaults to ${RELEASE}-${LIBC_FLAVOR}-${TAG}
variable "TAG" {
  default = "latest"
}

// SEMVER Semantic version to inject at build time
variable "SEMVER" {
  default = "undefined-docker"
}

// COMMITSHA Commit SHA to inject at build time
variable "COMMITSHA" {
  default = "undefined-docker"
}

// BUILDTIME Build time to inject at build time
variable "BUILDTIME" {
  default = "undefined-docker"
}

// TARGET The build target for the toolchain
variable "TARGET" {
  default = "x86_64-unknown-linux-gnu"
}

// BOOT_ARGS The boot arguments to bake into the UKI when generating an image
variable "BOOT_ARGS" {
  default = "debug console=ttyS0 rdinit=/usr/sbin/init"
}

// CROSSTOOLS_VER The version of crosstool-ng to use for building the toolchain
variable "CROSSTOOLS_VER" {
  default = "1.28.0"
}

// LINUX_VER The version of the linux kernel to use in this build
variable "LINUX_VER" {
  default = "6.19"
}

// BUSYBOX_VER The version of busybox to include with the initramfs
variable "BUSYBOX_VER" {
  default = "1.38.0"
}

// CORETUILS_VER The version of coreutils to use when building the rootfs
variable "COREUTILS_VER" {
  default = "0.9.0"
}

// CORETUILS_VER The version of uutils/shadow to use when building the rootfs
variable "SHADOW_VER" {
  default = "0.2.1"
}

// INIT_BUILDFLAGS The flags to pass to the Go compiler when building go-initd
variable "INIT_BUILDFLAGS" {
  default = "-v"
}

// INIT_LDFLAGS The flags to pass to the Go linker when building go-initd
variable "INIT_LDFLAGS" {
  default = "-s -w"
}

// SERVICED_BUILDFLAGS The flags to pass to the Go compiler when building go-serviced
variable "SERVICED_BUILDFLAGS" {
  default = "-v"
}

// SERVICED_INITFLAGS The flags to pass to the Go linker when building go-serviced
variable "SERVICED_INITFLAGS" {
  default = "-s -w"
}

target "_common" {
  platforms = ["linux/amd64"]
  output = ["type=docker"]
  labels = {
    "org.opencontainers.image.authors" = "Steven A. Zaluk <szaluk@stevenzaluk.com>"
    "org.opencontainers.image.source" = "https://github.com/stevezaluk-org/apathy-os"
    "org.opencontainers.image.url" = "https://github.com/stevezaluk-org/apathy-os"
    "org.opencontainers.image.version" = TAG
  }
}

target "toolchain" {
  inherits = ["_common"]

  dockerfile = "Dockerfile"
  context = "./build/toolchain"
  tags = ["ghcr.io/apathy-os/toolchain:${TAG}"]
  labels = {
    "org.opencontainers.title" = "cloud-host - toolchain"
  }

  args = {
    TARGET = TARGET
    CROSSTOOLS_VER = CROSSTOOLS_VER
  }
}

target "kernel" {
  inherits = ["_common"]

  dockerfile = "Dockerfile"
  context = "./build/kernel"
  tags = ["ghcr.io/apathy-os/kernel:${TAG}"]
  labels = {
    "org.opencontainers.title" = "cloud-host - kernel"
  }

  args = {
    LINUX_VER = LINUX_VER
  }
}

target "initramfs" {
  inherits = ["_common"]

  dockerfile = "Dockerfile"
  context = "./build/initramfs"

  contexts = {
    "initd" = "target:go-initd"
  }

  tags = ["ghcr.io/apathy-os/initramfs:${TAG}"]
  labels = {
    "org.opencontainers.title" = "cloud-host - initramfs"
  }

  args = {
    BUSYBOX_VER = BUSYBOX_VER
  }
}

target "rootfs" {
  inherits = ["_common"]

  dockerfile = "Dockerfile"
  context = "./build/rootfs"

  # We override locally here since we want to build "from scratch" when invoking through buildx bake
  contexts = {
    "toolchain" = "target:toolchain"
    "kernel" = "target:kernel"
    "initd" = "target:go-initd"
    "initramfs" = "target:initramfs"
  }

  tags = ["ghcr.io/apathy-os/rootfs:${TAG}"]
  labels = {
    "org.opencontainers.title" = "cloud-host - rootfs"
  }

  args = {
    BUSYBOX_VER = BUSYBOX_VER
    COREUTILS_VER = COREUTILS_VER
    SHADOW_VER = SHADOW_VER
  }
}

target "image" {
  inherits = ["_common"]

  dockerfile = "Dockerfile"
  context = "./build/image"
  contexts = {
    "kernel" = "target:kernel"
    "rootfs" = "target:rootfs"
    "initramfs" = "target:initramfs"
  }

  output = ["type=local,dest=out/apathy-os/image/"]

  args = {
    BOOT_ARGS = BOOT_ARGS
  }
}

target "go-initd" {
  inherits = ["_common"]

  dockerfile = "./build/initd/Dockerfile"
  context = "."
  tags = ["ghcr.io/apathy-os/go-initd:${TAG}"]

  labels = {
    "org.opencontainers.title" = "go-initd"
  }

  args = {
    INIT_BUILDFLAGS = INIT_BUILDFLAGS
    INIT_LDFLAGS = INIT_LDFLAGS
    SEMVER = SEMVER
    COMMITSHA = COMMITSHA
    BUILDTIME = BUILDTIME
  }
}

target "go-serviced" {
  inherits = ["_common"]

  dockerfile = "./build/serviced/Dockerfile"
  context = "."
  tags = ["ghcr.io/apathy-os/go-serviced:${TAG}"]

  labels = {
    "org.opencontainers.title" = "go-serviced"
  }

  args = {
    SERVICED_BUILDFLAGS = SERVICED_BUILDFLAGS
    SERVICED_INITFLAGS = SERVICED_INITFLAGS
    SEMVER = SEMVER
    COMMITSHA = COMMITSHA
    BUILDTIME = BUILDTIME
  }
}
