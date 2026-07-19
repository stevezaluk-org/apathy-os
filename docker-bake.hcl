// TAG The tag that gets assigned to a built image. Defaults to ${RELEASE}-${LIBC_FLAVOR}-${TAG}
variable "TAG" {
  default = "latest"
}

// TARGET The build target for the toolchain
variable "TARGET" {
  default = "x86_64-unknown-linux-gnu"
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
  default = "1.21.1"
}

// CORETUILS_VER The version of coreutils to use when building the rootfs
variable "COREUTILS_VER" {
  default = "0.9.0"
}

// TOOLCHAIN_VER The toolchain version to use when building the rootfs
variable "TOOLCHAIN_VER" {
  default = "latest"
}

// INITRAMFS_VER The version of the initramfs to use when building the rootfs
variable "INITRAMFS_VER" {
  default = "latest"
}

// AGENT_VER The version of the go-agentd to use when building the rootfs
variable "AGENT_VER" {
  default = "latest"
}

// AGENT_BUILDFLAGS The flags to pass to the Go compiler when building go-agentd
variable "AGENT_BUILDFLAGS" {
  default = "-v"
}

// AGENT_LDFLAGS The flags to pass to the Go linker when building go-agentd
variable "AGENT_LDFLAGS" {
  default = "-s -w"
}

// INIT_VER The version of go-initd to use when building the rootfs
variable "INIT_VER" {
  default = "latest"
}

// INIT_BUILDFLAGS The flags to pass to the Go compiler when building go-initd
variable "INIT_BUILDFLAGS" {
  default = "-v"
}

// INIT_LDFLAGS The flags to pass to the Go linker when building go-initd
variable "INIT_LDFLAGS" {
  default = "-s -w"
}

target "toolchain" {
  dockerfile = "build/toolchain/Dockerfile"
  platforms = ["linux/amd64"] // as of the moment, only x86_64 is supported
  tags = ["ghcr.io/immutable-os/toolchain:${TAG}"]
  output = ["type=docker"]
  labels = {
    "org.opencontainers.title" = "immutable-os - toolchain"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    TARGET = TARGET
    CROSSTOOLS_VER = CROSSTOOLS_VER
  }
}

target "kernel" {
  dockerfile = "Dockerfile"
  context = "./build/kernel"
  platforms = ["linux/amd64"]
  tags = ["ghcr.io/immutable-os/kernel:${TAG}"]
  output = ["type=docker"]
  labels = {
    "org.opencontainers.title" = "immutable-os - kernel"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    LINUX_VER = LINUX_VER
  }
}

target "initramfs" {
  dockerfile = "Dockerfile"
  context = "./build/initramfs"

  # We override locally here since we want to build "from scratch" when invoking through buildx bake
  contexts = {
    "ghcr.io/immutable-os/go-initd:latest" = "docker-image://ghcr.io/immutable-os/go-initd:latest"
    "ghcr.io/immutable-os/go-agentd:latest" = "docker-image://ghcr.io/immutable-os/go-agentd:latest"
  }

  platforms = ["linux/amd64"]
  tags = ["ghcr.io/immutable-os/initramfs:${TAG}"]
  output = ["type=docker"]
  labels = {
    "org.opencontainers.title" = "immutable-os - initramfs"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    BUSYBOX_VER = BUSYBOX_VER
  }
}

target "rootfs" {
  dockerfile = "Dockerfile"
  context = "./build/rootfs"

  # We override locally here since we want to build "from scratch" when invoking through buildx bake
  contexts = {
    "ghcr.io/immutable-os/go-initd:latest" = "docker-image://ghcr.io/immutable-os/go-initd:latest"
    "ghcr.io/immutable-os/go-agentd:latest" = "docker-image://ghcr.io/immutable-os/go-agentd:latest"
    "ghcr.io/immutable-os/kernel:latest" = "docker-image://ghcr.io/immutable-os/kernel:latest"
    "ghcr.io/immutable-os/initramfs:latest" = "docker-image://ghcr.io/immutable-os/initramfs:latest"
    "ghcr.io/immutable-os/toolchain:latest" = "docker-image://ghcr.io/immutable-os/toolchain:latest"
  }

  platforms = ["linux/amd64"]
  tags = ["ghcr.io/immutable-os/rootfs:${TAG}"]
  output = ["type=docker"]
  labels = {
    "org.opencontainers.title" = "immutable-os - rootfs"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    BUSYBOX_VER = BUSYBOX_VER
    LINUX_VER = LINUX_VER
    INIT_VER = INIT_VER
    AGENT_VER = AGENT_VER
    INITRAMFS_VER = INITRAMFS_VER
    TOOLCHAIN_VER = TOOLCHAIN_VER
    COREUTILS_VER = COREUTILS_VER
  }
}

target "go-agentd" {
  dockerfile = "go-agentd/Dockerfile"
  context = "./src"
  platforms = ["linux/amd64"]
  tags = ["ghcr.io/immutable-os/go-agentd:${TAG}"]
  output = ["type=docker"]

  labels = {
    "org.opencontainers.title" = "go-agentd"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    AGENT_BUILDFLAGS = AGENT_BUILDFLAGS
    AGENT_LDFLAGS = AGENT_LDFLAGS
  }
}

target "go-initd" {
  dockerfile = "go-initd/Dockerfile"
  context = "./src"
  platforms = ["linux/amd64"]
  tags = ["ghcr.io/immutable-os/go-initd:${TAG}"]
  output = ["type=docker"]

  labels = {
    "org.opencontainers.title" = "go-initd"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/immutable-os/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    INIT_BUILDFLAGS = INIT_BUILDFLAGS
    INIT_LDFLAGS = INIT_LDFLAGS
  }
}