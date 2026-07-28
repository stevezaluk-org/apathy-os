ARG AGENT_VER="latest"
ARG INIT_VER="latest"
ARG INITRAMFS_VER="latest"
ARG TOOLCHAIN_VER="latest"
ARG LINUX_VER="6.19"
ARG COREUTILS_VER="0.9.0"
ARG BUSYBOX_VER="1.38.0"

# This is a workaround since buildkit does not support expanding variables
# in `--from` arguments when using COPY. Pretty ugly but its the only easy
# way of doing this.
#
# See here: https://github.com/moby/buildkit/pull/3551
FROM ghcr.io/immutable-os/toolchain:$TOOLCHAIN_VER AS toolchain-build
FROM ghcr.io/immutable-os/go-initd:$INIT_VER AS initd-build
FROM ghcr.io/immutable-os/go-agentd:$AGENT_VER AS agentd-build
FROM ghcr.io/immutable-os/kernel:$LINUX_VER AS kernel-build
FROM ghcr.io/immutable-os/initramfs:$INITRAMFS_VER AS initramfs-build
FROM busybox:$BUSYBOX_VER-musl AS busybox-build

# This needs to be re-done to support our toolchain
FROM rust:1.96.0 AS  uutils-build
ARG COREUTILS_VER

WORKDIR /build/coreutils

ADD https://github.com/uutils/coreutils/archive/refs/tags/$COREUTILS_VER.tar.gz coreutils.tar.gz

RUN tar -xvf coreutils.tar.gz -C . --strip-components=1 && rm -rf coreutils.tar.gz
RUN cargo install \
          --locked \
          --path . \
          --bins \
          --root /build/coreutils

FROM debian:stable-slim AS rootfs-build
ARG BUSYBOX_VER
ARG AGENT_VER
ARG INIT_VER
ARG INITRAMFS_VER
ARG LINUX_VER

WORKDIR /build

RUN bash -c 'mkdir -p /build/{usr/bin,usr/sbin,etc,sys,dev,proc,tmp,boot,var,mnt}'

COPY --from=toolchain-build / /build
COPY --from=uutils-build /build/coreutils/bin/* /build/usr/bin
COPY --from=initd-build /app/init /build/init

# this must stay until we include /usr/bin/bash
COPY --from=busybox-build /bin/busybox /build/usr/sbin/busybox
COPY ./usr/share/manifest.d/system.yaml /build/usr/share/manifest.d/system.yaml
COPY ./usr/share/service.d/*.yaml /build/usr/share/service.d/
COPY ./usr/lib/os-release /build/usr/lib/os-release

# Symlink any required system directories
RUN ln -s usr/lib/os-release etc/os-release

WORKDIR /build/usr/sbin
RUN ln -s busybox findfs
RUN ln -s busybox switch_root

RUN chmod +x /build/usr/sbin/busybox

# Symlinks to our coreutils
RUN bash -c 'for prog in $(./usr/bin/coreutils --list); do ln -s ./coreutils usr/bin/$prog; done'

FROM scratch
ARG LINUX_VER
ARG INITRAMFS_VER

COPY --from=rootfs-build /build /
COPY --from=kernel-build /boot/vmlinuz /boot/vmlinuz
COPY --from=initramfs-build /boot/initrd.cpio /boot/initrd.img

SHELL ["/usr/sbin/busybox", "sh"]

ENV PATH="/usr/sbin:/usr/bin:/bin:/sbin"
ENV HOME="/"

CMD ["/usr/sbin/busybox", "sh"]