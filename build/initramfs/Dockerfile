ARG BUSYBOX_VER="1.38.0"

FROM busybox:$BUSYBOX_VER-musl AS busybox-build
FROM debian:stable-slim AS initramfs-build

RUN apt-get update \
    && \
    apt-get install cpio build-essential tar gzip -y

WORKDIR /build
RUN bash -c 'mkdir -p /build/{usr/sbin,usr/share/manifest.d,newroot,dev,proc,sys}'

COPY --from=busybox-build /bin/busybox usr/sbin/busybox

# This creates a confusing initramfs environment as agent is renamed to init, however this is a hard requirement of switch_root
COPY --from=ghcr.io/immutable-os/go-initd:latest /app/init /build/init
COPY --from=ghcr.io/immutable-os/go-agentd:latest /app/agentd /build/usr/sbin/init
COPY ./usr/share/manifest.d/init.yaml /build/usr/share/manifest.d/init.yaml

RUN bash -c 'find . | cpio -o -H newc > initrd.cpio'

FROM scratch

COPY --from=initramfs-build /build/initrd.cpio /boot/initrd.cpio