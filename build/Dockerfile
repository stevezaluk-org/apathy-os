FROM debian:stable-slim AS kernel-build

ARG LINUX_VER="6.19"
ARG LINUX_BUILD="true"
ARG LINUX_CONF="./build/linuxconfig"
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
RUN ["make", "-j", "20"]

FROM golang:1.25.5-alpine AS init-build

ARG INIT_BUILDFLAGS="-v"
ARG INIT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /out

COPY ./init .

RUN go build -o init $INIT_BUILDFLAGS -ldflags="$INIT_LDFLAGS" && chmod +x init

FROM golang:1.25.5-alpine AS agent-build

ARG AGENT_BUILDFLAGS="-v"
ARG AGENT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /out

COPY ./agent .
RUN go build -o agent $AGENT_BUILDFLAGS -ldflags="$AGENT_LDFLAGS"

FROM debian:stable-slim AS rootfs-build

WORKDIR /out/rootfs

RUN apt-get update \
    && \
    apt-get install tar gzip

RUN bash -c 'mkdir -p /rootfs/{bin,sbin,usr/bin,etc,sys,dev,proc,tmp}'

COPY --from=init-build /out/init /out/rootfs/usr/sbin/init
COPY --from=agent-build /out/agent /out/rootfs/usr/bin/AGENT_BUILDFLAGS

WORKDIR /out
RUN tar -czvf rootfs.tar.gz /out/rootfs

FROM debian:stable-slim AS initramfs-build

ARG E2FSPROGS_VER="1.47.4"

RUN apt-get update && apt-get install cpio build-essential tar gzip -y

WORKDIR /out/build/e2fsprogs

ADD https://github.com/tytso/e2fsprogs/archive/refs/tags/v$E2FSPROGS_VER.tar.gz e2fsprogs.tar.gz
RUN tar -xvf e2fsprogs.tar.gz -C . --strip-components=1

RUN ./configure LDFLAGS="-static" CFLAGS="-static" \
    && \
    make

WORKDIR /out/initramfs

RUN bash -c 'mkdir -p /out/initramfs/{bin,dev,proc,sys}'
# this needs to be externalized to ensure that we can build without this binary present
COPY "./build/busybox" /out/initramfs/bin/busybox

RUN mv /out/build/e2fsprogs/misc/mke2fs /out/initramfs/bin

COPY --from=init-build /out/init /out/initramfs/init
COPY --from=rootfs-build /out/rootfs.tar.gz /out/initramfs

RUN bash -c 'find . | cpio -o -H newc > ../initramfs.cpio'

FROM scratch

COPY --from=kernel-build /out/arch/x86_64/boot/bzImage /boot/vmlinuz
COPY --from=initramfs-build /out/initramfs.cpio /initramfs/initramfs.cpio
COPY --from=rootfs-build /out/rootfs.tar.gz /rootfs/rootfs.tar.gz