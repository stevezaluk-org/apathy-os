ARG ROOTFS_VER="latest"
ARG BOOT_ARGS=""

FROM ghcr.io/immutable-os/rootfs:$ROOTFS_VER AS rootfs-build
FROM debian:stable-slim AS image-build

WORKDIR /build

RUN apt-get update && \
    apt-get install genimage systemd-ukify mtools dosfstools e2fsprogs -y

RUN mkdir -p /build/esp/EFI/BOOT

COPY ./release.cfg /build
COPY --from=rootfs-build / /rootfs

RUN ukify build \
      --linux=/rootfs/boot/vmlinuz \
      --initrd=/rootfs/boot/initrd.img \
      --cmdline=$BOOT_ARGS \
      --output=/build/esp/EFI/BOOT/BOOTX64.EFI

RUN	genimage \
      --config /build/release.cfg \
      --rootpath=/rootfs \
      --tmppath=/build/tmp \
      --inputpath=/build/esp \
      --outputpath=/build/image

FROM scratch

COPY --from=image-build /build/image/disk.img /