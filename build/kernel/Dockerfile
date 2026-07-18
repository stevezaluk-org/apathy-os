ARG LINUX_VER="6.19"

FROM debian:stable-slim AS kernel-build
ARG LINUX_VER

WORKDIR /build

RUN apt-get update \
    && \
    apt-get install bc bison flex libssl-dev libelf-dev build-essential pkg-config gcc make tar gzip rsync -y

ADD https://github.com/torvalds/linux/archive/refs/tags/v$LINUX_VER.tar.gz /build/linux.tar.gz

RUN tar -xvf linux.tar.gz -C . --strip-components=1 && rm linux.tar.gz

# Doing an entire build here probs wont make sense long term but it will do for now
COPY ./config.default ./.config

# merge old config with paramters introduced in newer versions
RUN make olddefconfig
RUN make -j $(nproc)

FROM scratch
ARG LINUX_VER

COPY --from=kernel-build /build/arch/x86_64/boot/bzImage /boot/vmlinuz-$LINUX_VER