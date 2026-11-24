ARG TOOLCHAIN_VER="2.42"
ARG RUST_VER="1.96.0"

FROM ghcr.io/stevezaluk-org/apathy-os/toolchain:$TOOLCHAIN_VER AS toolchain-build
FROM debian:trixie-slim@sha256:a72b4433de1cf7d6bae0c4cee74d0547f58601a643b3375f070c51600d63b6b2
ARG RUST_VER

# ca-certificates required as a pre-req for using the snapshot archive
RUN apt-get update -y && \
    apt-get install ca-certificates -y

# Remove any old sources that exist for APT so that we can use our snapshot
RUN rm -rf /etc/apt/sources.list.d/*
COPY etc/apt/sources.list.d/* /etc/apt/sources.list.d/

RUN apt-get update -y && \
    apt-get install gcc g++ pkg-config meson tar unzip rustup -y

# Install our rust compiler for other dependencies that depend on this
RUN rustup toolchain install $RUST_VER --profile=minimal && \
    rustup default $RUST_VER

# Always ensure that we have our custom toolchain built so that we can build our components against it
COPY --from=toolchain-build / /opt/sysroot

# Configuration files for: meson, cargo
COPY opt/sysroot/conf/meson-sysroot.ini /opt/conf/meson-sysroot.ini
COPY opt/sysroot/conf/cargo.toml /etc/cargo/config.toml

# environment variables required so that the C compiler picks up our sysroot
ENV CC="gcc --sysroot=/opt/sysroot"
ENV CXX="g++ --sysroot=/opt/sysroot"
ENV CPPFLAGS="-I/opt/sysroot/usr/include"
ENV LDFLAGS="-L/opt/sysroot/usr/lib -Wl,-rpath-link,/opt/sysroot/usr/lib"
ENV BINDGEN_EXTRA_CLANG_ARGS="--sysroot=/opt/sysroot"
ENV PKG_CONFIG_SYSROOT_DIR="/opt/sysroot"
ENV PKG_CONFIG_LIBDIR="/opt/sysroot/usr/lib/pkgconfig"
ENV CARGO_HOME="/etc/cargo"

SHELL ["/bin/bash", "-o", "pipefail", "-c"]
