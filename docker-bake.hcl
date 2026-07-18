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

// AGENT_BUILDFLAGS The flags to pass to the Go compiler when building go-agentd
variable "AGENT_BUILDFLAGS" {
  default = "-v"
}

// AGENT_LDFLAGS The flags to pass to the Go linker when building go-agentd
variable "AGENT_LDFLAGS" {
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