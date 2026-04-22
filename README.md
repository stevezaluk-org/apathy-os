# What is immutable-os?
Immutable-os is a fully immutable, Linux distribution built from scratch that allows you to build purpose-built full
operating systems for your use case. Immutable-os was not built on top of a pre-existing distro, does not depend on
existing Linux tooling (like systemd) and does not have a shell

## Source Tree
* ``build`` - Config files, docker/bake files for constructing system components and building tooling
* ``initd`` - PID 1 for immutable-os. Handles mounting the rootfs, device initialization, and handling the graceful shutdown process

## Building
Since immutable-os is designed as a "purpose-built" operating system, there is no unified ISO installer like a more mainstream
distribution like Ubuntu would provide. As a result, you will need to compile each of the system components and build
a disk image from scratch.

Immutable-os uses a containerized build process to ensure that builds are modular, reproducible, and fully contained. To
build immutable-os, please see the below dependencies:

### Dependencies
The below are the minimum requirements for building immutable-os:
* Any linux distribution
* Docker (with buildx/bake support)
* make
* wget
* ukify
* genimage

Although not required, if you are doing development work on immutable os, the following is recommended in addition to
the above:
* go 1.26.1 compiler - Required if you plan on modifying, developing, or testing any system components
* OVMF firmware - Required for virtualizing UEFI-based systems. Needs to be installed to ``/usr/share/OVMF``
* qemu - Required if you plan on virtualizing an x86-64 machine running immutable-os

### Variants
immutable-os offers two primary variants: release and development. The names of each of these variants are pretty self-explanatory,
however one important consideration to be aware of is that the development release will drop you into a busybox shell
after the root filesystem is mounted. This is provided to ensure that you have a way of properly testing new features
that are implemented. Otherwise, the two releases are identical. This may change in the future.

In addition to the variants provided above, you have the ability to select the flavor of libc that you would like to
bundle with your installation. There are three flavors available:
* static (default) - No libc implementation
* glibc - Builds/installs the standard GNU libc implementation
* musl - Builds/installs the musl libc implementation

### Build
Building immutable-os is very simple:

```
# Download any required dependencies (just the x86_64 variant of busybox for now)
make deps

# Build rootfs, initramfs, initd, and any other required system components
make build RELEASE=<variant> LIBC_FLAVOR=<flavor>  

# Build a UEFI disk image file that can either be virtualized with qemu or flashed to a physical drive
make image

# (Optional) Clean up any build artifacts required for building the UEFI image
make clean
```
