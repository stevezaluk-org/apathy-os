ARG RELEASE="release"
ARG LIBC_FLAVOR="static"

FROM debian:stable-slim AS kernel-build

ARG LINUX_VER="6.19"
ARG LINUX_BUILD="true"
ARG LINUX_CONF="./build/kernel/config.default"
ARG CORES=$(nproc)

WORKDIR /out

RUN apt-get update \
    && \
    apt-get install bc bison flex libssl-dev libelf-dev build-essential pkg-config gcc make tar gzip -y

ADD https://github.com/torvalds/linux/archive/refs/tags/v$LINUX_VER.tar.gz /out/linux.tar.gz

RUN tar -xvf linux.tar.gz -C . --strip-components=1 \
    && \
    rm linux.tar.gz

# Doing an entire build here probs wont make sense long term but it will do for now
COPY $LINUX_CONF ./.config

# merge old config with paramters introduced in newer versions
RUN make olddefconfig
RUN make -j $(nproc)

FROM golang:1.25.5-alpine AS init-build

ARG INIT_BUILDFLAGS="-v"
ARG INIT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /out

COPY ./init .

RUN go build -o init $INIT_BUILDFLAGS -ldflags="$INIT_LDFLAGS" && chmod +x init

FROM debian:stable-slim AS rootfs-glibc

ARG GLIBC_VER="2.43"

WORKDIR /out

RUN bash -c 'mkdir -p /out/rootfs/{bin,sbin,usr/bin,etc,sys,dev,proc,tmp}'

WORKDIR /out/build/glibc

RUN apt-get update \
    && \
    apt-get install tar gzip gcc make gawk bison python3 -y

ADD https://ftp.gnu.org/gnu/glibc/glibc-2.43.tar.gz glibc.tar.gz

RUN mkdir build && tar -xvf glibc.tar.gz -C build/ --strip-components=1

RUN ./build/configure --prefix=/out/rootfs/usr --datadir=$(pwd)
RUN make -j $(nproc)
RUN make install

FROM debian:stable-slim AS rootfs-static

WORKDIR /out

# not much really needed here as the whole point of a static build is not needing a libc implementation
RUN bash -c 'mkdir -p /out/rootfs/{bin,sbin,usr/bin,etc,sys,dev,proc,tmp}'

FROM rootfs-${LIBC_FLAVOR} AS rootfs-build

WORKDIR /out

COPY "./build/deps/busybox-x86_64" /out/rootfs/bin/busybox
COPY --from=init-build /out/init /out/rootfs/usr/sbin/init

RUN tar -czvf rootfs.tar.gz /out/rootfs

FROM debian:stable-slim AS initramfs-build

RUN apt-get update \
    && \
    apt-get install cpio build-essential tar gzip -y

WORKDIR /out/initramfs

RUN bash -c 'mkdir -p /out/initramfs/{bin,dev,proc,sys}'

COPY "./build/deps/busybox-x86_64" /out/initramfs/bin/busybox

COPY --from=init-build /out/init /out/initramfs/init

RUN bash -c 'find . | cpio -o -H newc > ../initramfs.cpio'

FROM scratch

COPY --from=kernel-build /out/arch/x86_64/boot/bzImage /boot/vmlinuz
COPY --from=initramfs-build /out/initramfs.cpio /initramfs/initramfs.cpio
COPY --from=rootfs-build /out/rootfs.tar.gz /rootfs/rootfs.tar.gz