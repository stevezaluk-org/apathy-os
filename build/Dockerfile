ARG RELEASE="release"

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

RUN unset LD_AUDIT \
    unset LD_PRELOAD

WORKDIR /out

RUN bash -c 'mkdir -p /out/rootfs/{usr/bin,usr/sbin,etc,sys,dev,proc,tmp,boot,var,mnt}'

COPY --from=kernel-build /out/headers/include /out/rootfs/usr/include

WORKDIR /out/build/glibc

RUN apt-get update \
    && \
    apt-get install tar gzip gcc g++ make gawk bison python3 libgmp-dev libmpfr-dev libmpc-dev libisl-dev zlib1g-dev libc6-dev -y

ADD https://ftp.gnu.org/gnu/glibc/glibc-$GLIBC_VER.tar.gz glibc.tar.gz

RUN mkdir build && tar -xvf glibc.tar.gz -C build/ --strip-components=1

WORKDIR /out/build/glibc
RUN ./build/configure --prefix=/usr --datadir=$(pwd) --with-sysroot=/out/rootfs
RUN make -j $(nproc) slibdir=/usr/lib64 rtlddir=/usr/lib64
RUN make DESTDIR=/out/rootfs install slibdir=/usr/lib64 rtlddir=/usr/lib64 sbindir=/usr/sbin rootsbindir=/usr/sbin

WORKDIR /out/build/gcc

ADD https://ftp.gnu.org/gnu/gcc/gcc-${GCC_VER}/gcc-${GCC_VER}.tar.gz gcc.tar.gz

RUN mkdir build && tar -xvf gcc.tar.gz -C build/ --strip-components=1 && rm -rf gcc.tar.gz

RUN ./build/configure \
    --with-sysroot=/out/rootfs \
    --prefix=/usr \
    --enable-languages=c,c++ \
    --enable-shared=libgcc,libstdc++ \
    --disable-bootstrap \
    --disable-multilib

# [WARN] This is a hack and should really be rectified by either building both 32-bit and 64-bit stl's or just 64-bit
RUN mkdir -p /out/rootfs/usr/lib && cp -r /out/rootfs/usr/lib64/* /out/rootfs/usr/lib

# We only really care about the gcc shared libs and not the compiler itself so we only target that here
RUN make -j$(nproc) all-gcc
RUN make -j$(nproc) all-target-libgcc
RUN make -j$(nproc) all-target-libstdc++-v3

RUN make DESTDIR=/out/rootfs install-target-libgcc
RUN make DESTDIR=/out/rootfs install-target-libstdc++-v3

WORKDIR /out

COPY "./build/manifest.d/system.yaml" /out/rootfs/usr/share/manifest.d/system.yaml
COPY "./build/service.d/*.yaml" /out/rootfs/usr/share/service.d/

# this must stay until we include /usr/bin/bash
COPY "./build/deps/busybox-x86_64" /out/rootfs/usr/sbin/busybox
COPY --from=uutils-build /build/coreutils/bin/* /out/rootfs/usr/bin
COPY --from=init-build /out/init /out/rootfs/usr/sbin/init
COPY --from=agentd-build /out/agent /out/rootfs/init

RUN rm -rf /out/rootfs/out

WORKDIR /out/rootfs

# Symlinks for our core directories
RUN ln -s usr/bin bin
RUN ln -s usr/sbin sbin
RUN ln -s usr/lib lib
RUN ln -s usr/lib64 lib64

# Symlinks to our coreutils
RUN bash -c 'for prog in $(./usr/bin/coreutils --list); do ln -s ./coreutils usr/bin/$prog; done'

FROM scratch

COPY --from=rootfs-build /out/rootfs /
COPY --from=kernel-build /out/arch/x86_64/boot/bzImage /boot/vmlinuz
COPY --from=initramfs-build /out/initrd.cpio /boot/initrd.img