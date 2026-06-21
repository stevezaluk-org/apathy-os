ARG RELEASE="release"

FROM debian:stable-slim AS kernel-build
ARG LINUX_VER="6.19"

WORKDIR /out

RUN apt-get update \
    && \
    apt-get install bc bison flex libssl-dev libelf-dev build-essential pkg-config gcc make tar gzip -y

ADD https://github.com/torvalds/linux/archive/refs/tags/v$LINUX_VER.tar.gz /out/linux.tar.gz

RUN tar -xvf linux.tar.gz -C . --strip-components=1 \
    && \
    rm linux.tar.gz

# Doing an entire build here probs wont make sense long term but it will do for now
COPY "./build/kernel/config.default" ./.config

# merge old config with paramters introduced in newer versions
RUN make olddefconfig
RUN make -j $(nproc)

FROM golang:1.26.1-alpine AS init-build

ARG INIT_BUILDFLAGS="-v"
ARG INIT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /out

COPY ./initd initd/
COPY ./go-manifest go-manifest/

# this can be removed once this goes public
RUN go work init ./initd ./go-manifest

RUN go build -o init $INIT_BUILDFLAGS -ldflags="$INIT_LDFLAGS" ./initd && chmod +x init

FROM golang:1.26.1-alpine AS agentd-build

ARG INIT_BUILDFLAGS="-v"
ARG INIT_LDFLAGS="-s -w"

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /out

COPY ./agentd agentd/
COPY ./go-manifest go-manifest/

# this can be removed once this goes public
RUN go work init ./agentd ./go-manifest

RUN go build -o agent $INIT_BUILDFLAGS -ldflags="$INIT_LDFLAGS" ./agentd && chmod +x agent

FROM rust:1.96.0 AS  uutils-build
ARG COREUTILS_VER="0.9.0"

WORKDIR /build/coreutils

ADD https://github.com/uutils/coreutils/archive/refs/tags/$COREUTILS_VER.tar.gz coreutils.tar.gz

RUN tar -xvf coreutils.tar.gz -C . --strip-components=1 && rm -rf coreutils.tar.gz
RUN cargo install \
          --locked \
          --path . \
          --bins \
          --root /build/coreutils

FROM debian:stable-slim AS rootfs-build
ARG GLIBC_VER="2.43"
ARG GCC_VER="16.1.0"

WORKDIR /out

RUN bash -c 'mkdir -p /out/rootfs/{bin,sbin,usr/bin,etc,sys,dev,proc,tmp,boot,var,mnt}'

WORKDIR /out/build/glibc

RUN apt-get update \
    && \
    apt-get install tar gzip gcc g++ make gawk bison python3 libgmp-dev libmpfr-dev libmpc-dev libisl-dev zlib1g-dev -y

ADD https://ftp.gnu.org/gnu/glibc/glibc-$GLIBC_VER.tar.gz glibc.tar.gz

RUN mkdir build && tar -xvf glibc.tar.gz -C build/ --strip-components=1

WORKDIR /out/build/glibc
RUN ./build/configure --prefix=/out/rootfs/usr --datadir=$(pwd)
RUN make -j $(nproc)
RUN make install

WORKDIR /out/build/gcc

ADD https://ftp.gnu.org/gnu/gcc/gcc-${GCC_VER}/gcc-${GCC_VER}.tar.gz gcc.tar.gz

RUN mkdir build && tar -xvf gcc.tar.gz -C build/ --strip-components=1 && rm -rf gcc.tar.gz

RUN ./build/configure \
     --prefix=/out/rootfs/usr \
    --enable-languages=c,c++ \
    --enable-shared=libgcc,libstdc++ \
    --disable-bootstrap \
    --disable-multilib

# We only really care about the gcc shared libs and not the compiler itself so we only target that here
RUN make all-gcc -j $(nproc)
RUN make all-target-libgcc  -j $(nproc)
RUN make all-target-libstdc++-v3 -j $(nproc)

RUN make install-target-libgcc
RUN make install-target-libstdc++-v3

WORKDIR /out

COPY "./build/manifest.d/system.yaml" /out/rootfs/usr/share/manifest.d/system.yaml
COPY "./build/service.d/*.yaml" /out/rootfs/usr/share/service.d/

COPY "./build/deps/busybox-x86_64" /out/rootfs/usr/sbin/busybox
COPY --from=uutils-build /build/coreutils/bin/* /out/rootfs/usr/bin
COPY --from=init-build /out/init /out/rootfs/usr/sbin/init
COPY --from=agentd-build /out/agent /out/rootfs/init

WORKDIR /out/rootfs

RUN ln -s usr/lib /out/rootfs/lib
RUN ln -s usr/lib /out/rootfs/lib64

FROM debian:stable-slim AS initramfs-build

RUN apt-get update \
    && \
    apt-get install cpio build-essential tar gzip -y

WORKDIR /out/initramfs

RUN bash -c 'mkdir -p /out/initramfs/{usr/sbin,usr/share/manifest.d,newroot,dev,proc,sys}'
RUN bash -c 'mkdir -p /out/initramfs/newroot/{etc,var,mnt}'

COPY "./build/deps/busybox-x86_64" /out/initramfs/usr/sbin/busybox
COPY --from=init-build /out/init /out/initramfs/usr/sbin/init

# This creates a confusing initramfs environment as agent is renamed to init, however this is a hard requirement of switch_root
COPY --from=agentd-build /out/agent /out/initramfs/init

COPY "./build/manifest.d/init.yaml" /out/initramfs/usr/share/manifest.d/init.yaml

RUN bash -c 'find . | cpio -o -H newc > ../initrd.cpio'

FROM scratch

COPY --from=rootfs-build /out/rootfs /
COPY --from=kernel-build /out/arch/x86_64/boot/bzImage /boot/vmlinuz
COPY --from=initramfs-build /out/initrd.cpio /boot/initrd.img