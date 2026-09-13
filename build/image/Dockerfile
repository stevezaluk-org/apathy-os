ARG ROOTFS_VER="latest"
ARG INITRAMFS_VER="latest"
ARG LINUX_VER="latest"
ARG BOOT_ARGS="debug root=87965c58-36e0-11f1-8de9-0242ac120002 console=ttyS0 rdinit=/usr/sbin/init"

FROM ghcr.io/immutable-os/rootfs:$ROOTFS_VER AS rootfs-build
FROM ghcr.io/immutable-os/initramfs:$INITRAMFS_VER AS initramfs-build
FROM ghcr.io/immutable-os/kernel:$LINUX_VER AS kernel-build
FROM debian:stable-slim AS image-build
ARG BOOT_ARGS

RUN apt-get update && \
    apt-get install cpio genimage systemd-ukify mtools dosfstools e2fsprogs -y

COPY --from=rootfs-build / /rootfs
COPY --from=initramfs-build / /initramfs
COPY --from=kernel-build /boot/vmlinuz /kernel/vmlinuz

WORKDIR /initramfs

RUN bash -c 'find . | cpio -o -H newc > initrd.img'

WORKDIR /build

RUN mkdir -p /build/esp/EFI/BOOT

RUN ukify build \
      --linux=/kernel/vmlinuz \
      --initrd=/initramfs/initrd.img \
      --cmdline="$BOOT_ARGS" \
      --os-release=/rootfs/usr/lib/os-release \
      --output=/build/esp/EFI/BOOT/BOOTX64.EFI

COPY ./release.cfg /build
RUN	genimage \
      --config /build/release.cfg \
      --rootpath=/rootfs \
      --tmppath=/build/tmp \
      --inputpath=/build/esp \
      --outputpath=/build/image

FROM scratch

COPY --from=image-build /build/image/disk.img /