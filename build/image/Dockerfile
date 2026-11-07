ARG BOOT_ARGS="debug console=ttyS0 rdinit=/usr/sbin/init"

FROM debian:stable-slim AS image-build
ARG BOOT_ARGS

RUN apt-get update && \
    apt-get install cpio genimage systemd-ukify mtools dosfstools e2fsprogs squashfs-tools systemd-boot-efi -y

RUN mkdir -p /build/esp/EFI/BOOT && \
    mkdir -p /build/esp/EFI/Linux

RUN cp /usr/lib/systemd/boot/efi/systemd-bootx64.efi /build/esp/EFI/BOOT/BOOTX64.EFI
COPY ./boot/loader/loader.conf /build/esp/loader/loader.conf
COPY ./release.cfg /build

COPY --from=rootfs / /rootfs
COPY --from=initramfs / /initramfs
COPY --from=kernel /boot/vmlinuz /kernel/vmlinuz

WORKDIR /initramfs

RUN mksquashfs /rootfs /initramfs/rootfs.squashfs -comp lz4 -reproducible -all-root -noappend
RUN bash -c 'find . -mindepth 1 | sort | cpio -o -H newc -R 0:0 --quiet > /build/initrd.img'

WORKDIR /build

RUN ukify build \
      --linux=/kernel/vmlinuz \
      --initrd=/build/initrd.img \
      --cmdline="$BOOT_ARGS" \
      --os-release=@/rootfs/usr/lib/os-release \
      --stub=/usr/lib/systemd/boot/efi/linuxx64.efi.stub \
      --output=/build/esp/EFI/Linux/cloud-host_A.efi

RUN	genimage \
      --config /build/release.cfg \
      --rootpath=/rootfs \
      --tmppath=/build/tmp \
      --inputpath=/build/esp \
      --outputpath=/build/image

FROM scratch

COPY --from=image-build /build/image/disk.img /
