FROM debian:stable-slim AS initramfs-build
ARG BUSYBOX_VER="1.21.1"

RUN apt-get update \
    && \
    apt-get install cpio build-essential tar gzip -y

WORKDIR /build
RUN bash -c 'mkdir -p /build/{usr/sbin,usr/share/manifest.d,newroot,dev,proc,sys}'

ADD https://busybox.net/downloads/binaries/$BUSYBOX_VER/busybox-x86_64 usr/sbin/busybox

# This creates a confusing initramfs environment as agent is renamed to init, however this is a hard requirement of switch_root
COPY --from=ghcr.io/immutable-os/go-initd:latest /app/initd /build/init
COPY --from=ghcr.io/immutable-os/go-agentd:latest /app/agentd /build/usr/sbin/init
COPY ./usr/share/manifest.d/init.yaml /build/usr/share/manifest.d/init.yaml

RUN bash -c 'find . | cpio -o -H newc > initrd.cpio'

FROM scratch

COPY --from=initramfs-build /build/initrd.cpio /boot/initrd.cpio