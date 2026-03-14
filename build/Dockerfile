FROM debian:stable-slim AS kernel-build

ARG LINUX_VER="6.19"
ARG LINUX_BUILD="true"
ARG LINUX_CONF="./build/linuxconfig"
ARG CORES=$(nproc)

WORKDIR /kern

RUN apt-get update \
    && \
    apt-get install bc bison flex libssl-dev libelf-dev build-essential pkg-config gcc make tar gzip -y

ADD https://github.com/torvalds/linux/archive/refs/tags/v$LINUX_VER.tar.gz /kern/linux.tar.gz

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

WORKDIR /build

COPY ./init .

RUN go build -o init $INIT_BUILDFLAGS -ldflags="$INIT_LDFLAGS"

FROM debian:stable-slim AS initramfs-build

WORKDIR /initramfs

RUN bash -c 'mkdir -p /initramfs/{bin,dev,proc,sys}'

COPY --from=init-build /build/init /initramfs/bin/init
RUN ln -s /initramfs/bin/init /init

FROM scratch

COPY --from=kernel-build /kern/arch/x86_64/boot/bzImage /boot/vmlinuz
COPY --from=initramfs-build /initramfs /initramfs