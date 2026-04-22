// RELEASE Allows you to define the release of container-os you would like to build. A development release includes a shell
variable "RELEASE" {
  default = "release"
}

// LIBC_FLAVOR If your use case requires dynamic linking then you can select a libc flavor to install. Static is the default
variable "LIBC_FLAVOR" {
  default = "static"
}

// TAG The tag that gets assigned to a built image. Defaults to ${RELEASE}-${LIBC_FLAVOR}-${TAG}
variable "TAG" {
  default = "latest"
}

// GLIBC_VER The version of glibc that you want to use when building glibc flavors
variable "GLIBC_VER" {
  default = "2.43"
}

// MUSL_VER The version of musl that you want to use when building musl flavors
variable "MUSL_VER" {
  default = "1.2.6"
}

// LINUX_VER The version of of the linux kernel that should be built
variable "LINUX_VER" {
  default = ""
}


// _common The base for all builds
target "_common" {
  dockerfile = "build/Dockerfile"
  platforms = ["linux/amd64"] // as of the moment, only x86_64 is supported
  tags = ["ghcr.io/stevezaluk/container-os/container-os:${RELEASE}-${LIBC_FLAVOR}-${TAG}"]
  output = ["type=local,dest=dist/out/${RELEASE}-${LIBC_FLAVOR}"]
  labels = {
    "org.opencontainers.title" = "immutable-os (${RELEASE}/${LIBC_FLAVOR})"
    "org.opencontainers.image.authors" = "Steven A. Zaluk <stevezaluk@gmail.com>"
    "org.opencontainers.image.source" = "https://github.com/stevezaluk/immutable-os"
    "org.opencontainers.image.url" = "https://github.com/stevezaluk/immutable-os"
    "org.opencontainers.image.version" = TAG
  }

  args = {
    RELEASE = RELEASE
    LIBC_FLAVOR = LIBC_FLAVOR
    LINUX_VER = LINUX_VER
  }
}

group "matrix" {
  targets = [
    "release-glibc",
    "development-glibc",
    "release-musl",
    "development-musl",
    "release-static",
    "development-static"
  ]
}

target "release-glibc" {
  inherits = ["_common"]
  output = ["type=local,dest=dist/out/release-glibc"]
  args = {
    RELEASE = "release"
    LIBC_FLAVOR = "glibc"
    GLIBC_VER = GLIBC_VER
  }
}

target "development-glibc" {
  inherits = ["_common"]
  output = ["type=local,dest=dist/out/development-glibc"]
  args = {
    RELEASE = "development"
    LIBC_FLAVOR = "glibc"
    GLIBC_VER = GLIBC_VER
  }
}

target "release-musl" {
  inherits = ["_common"]
  output = ["type=local,dest=dist/out/release-musl"]
  args = {
    RELEASE = "release"
    LIBC_FLAVOR = "musl"
    MUSL_VER = MUSL_VER
  }
}

target "development-musl" {
  inherits = ["_common"]
  output = ["type=local,dest=dist/out/development-musl"]
  args = {
    RELEASE = "development"
    LIBC_FLAVOR = "musl"
    MUSL_VER = MUSL_VER
  }
}

target "release-static" {
  inherits = ["_common"]
  output = ["type=local,dest=dist/out/release-static"]
  args = {
    RELEASE = "release"
    LIBC_FLAVOR = "static"
  }
}

target "development-glibc" {
  inherits = ["_common"]
  output = ["type=local,dest=dist/out/development-static"]
  args = {
    RELEASE = "development"
    LIBC_FLAVOR = "static"
  }
}
