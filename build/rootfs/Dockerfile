ARG TOOLCHAIN_VER="2.42"
ARG COREUTILS_VER="0.9.0"
ARG BUSYBOX_VER="1.38.0"
ARG UTILLINUX_VER="2.42.2"
ARG SHADOW_VER="0.2.0"
ARG LIBPAM_VER="1.7.2"
ARG LIBCRYPT_VER="4.5.2"

# This is a workaround since buildkit does not support expanding variables
# in `--from` arguments when using COPY. Pretty ugly but its the only easy
# way of doing this.
#
# See here: https://github.com/moby/buildkit/pull/3551
FROM busybox:$BUSYBOX_VER-musl AS busybox-build
FROM ghcr.io/stevezaluk-org/apathy-os/toolchain:$TOOLCHAIN_VER AS toolchain-build
FROM ghcr.io/stevezaluk-org/apathy-os/uutils-coreutils:$COREUTILS_VER AS uutils-build
FROM ghcr.io/stevezaluk-org/apathy-os/uutils-shadow:$SHADOW_VER AS shadow-build
FROM ghcr.io/stevezaluk-org/apathy-os/util-linux:$UTILLINUX_VER AS utillinux-build
FROM ghcr.io/stevezaluk-org/apathy-os/libcrypt:$LIBCRYPT_VER AS libcrypt-build
FROM ghcr.io/stevezaluk-org/apathy-os/libpam:$LIBPAM_VER AS libpam-build

FROM ghcr.io/stevezaluk-org/apathy-os/base:latest AS rootfs-build
ARG BUSYBOX_VER

WORKDIR /build

RUN bash -c 'mkdir -p /build/{usr/bin,usr/sbin,etc,sys,dev,proc,tmp,boot,var,mnt}'

# The toolchain makes up the base of our root filesystem
COPY --from=toolchain-build / /build

# libcrypt/libpam are both required at runtime for user management functions (no headers)
COPY --from=libcrypt-build /usr/lib/* /build/usr/lib
COPY --from=libpam-build /usr/lib/x86_64-linux-gnu/ /build/usr/lib
COPY --from=libpam-build /usr/lib/security /build/usr/lib/security
COPY --from=libpam-build /usr/share/pam /build/usr/share
COPY --from=libpam-build /usr/sbin /build/usr/sbin

# uutils-coreutils, shadow, and util-linux makes up majority of the user space binaries
COPY --from=uutils-build /usr/bin/coreutils /build/usr/bin/coreutils
COPY --from=shadow-build /usr/bin/* /build/usr/bin
COPY --from=utillinux-build /usr/bin/* /build/usr/bin
COPY --from=utillinux-build /usr/sbin/* /build/usr/sbin
COPY --from=utillinux-build /usr/lib/* /build/usr/lib

# Then our custom components
COPY --from=initd /app/init /build/init

# this must stay until we include /usr/bin/bash
COPY --from=busybox-build /bin/busybox /build/usr/sbin/busybox
RUN cp /build/usr/sbin/busybox /build/usr/bin/sh

COPY ./usr/share/manifest.d/system.yaml /build/usr/share/manifest.d/system.yaml
COPY ./usr/share/service.d/*.json /build/usr/share/service.d/
COPY ./usr/lib/os-release /build/usr/lib/os-release

# Symlink any required system directories
RUN ln -s usr/lib/os-release etc/os-release
RUN chmod +x /build/usr/sbin/busybox

WORKDIR /build/etc

# Required template files for user control. User control is generally handled by go-manifest, but we always need a root
# user so the template contains entries for this
COPY ./etc/passwd .
COPY ./etc/group .
COPY ./etc/shadow .
COPY ./etc/issue .
COPY ./etc/login.defs .

# All configuration files here are required for configuring PAM
COPY ./etc/pam.d /build/etc/pam.d

WORKDIR /build

# Symlinks to our coreutils
RUN bash -c 'for prog in $(./usr/bin/coreutils --list); do ln -s ./coreutils usr/bin/$prog; done'

FROM scratch
COPY --from=rootfs-build /build /

ENV PATH="/usr/sbin:/usr/bin:/bin:/sbin"
ENV HOME="/"

SHELL ["/usr/sbin/busybox", "sh"]
CMD ["/usr/sbin/busybox", "sh"]
