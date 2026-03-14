ARG LINUX_VER="6.19"
ARG LINUX_BUILD="true"
ARG LINUX_CONF="./linuxconfig"
ARG CORES=$(nproc)

FROM debian:stable-slim AS kernel-build

RUN mkdir -p /kern
WORKDIR /kern

RUN apt-get update \
    && \
    apt-get install wget git bc bison flex libssl-dev libelf-dev build-essential pkg-config gcc make tar gzip -y

RUN wget https://github.com/torvalds/linux/archive/refs/tags/v$LINUX_VER.tar.gz -O linux.tar.gz \
    && \
    tar -xvf linux.tar.gz -C . --strip-components=1 \
    && \
    rm linux.tar.gz

# Doing an entire build here probs wont make sense long term but it will do for now
COPY $LINUX_CONF ./.config

# merge old config with paramters introduced in newer versions
RUN make olddefconfig
RUN ["make", "-j", "20"]

FROM scratch

COPY --from=kernel-build /kern/arch/x86_64/boot/bzImage /boot/vmlinuz