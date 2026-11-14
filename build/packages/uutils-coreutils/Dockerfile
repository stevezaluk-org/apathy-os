ARG LIBPAM_VER="1.7.2"
ARG LIBCRYPT_VER="4.5.2"
ARG COREUTILS_VER="0.9.0"

FROM ghcr.io/stevezaluk-org/apathy-os/libcrypt:$LIBCRYPT_VER AS libcrypt-build
FROM ghcr.io/stevezaluk-org/apathy-os/libpam:$LIBPAM_VER AS libpam-build
FROM ghcr.io/stevezaluk-org/apathy-os/base:latest AS builder
ARG LIBPAM_VER
ARG LIBCRYPT_VER
ARG COREUTILS_VER

WORKDIR /build/coreutils

ADD https://github.com/uutils/coreutils/archive/refs/tags/$COREUTILS_VER.tar.gz coreutils.tar.gz

COPY --from=libcrypt-build /usr/lib/* /opt/sysroot/usr/lib
COPY --from=libcrypt-build /usr/include /opt/sysroot/usr/include
COPY --from=libpam-build /usr/lib/x86_64-linux-gnu/ /opt/sysroot/usr/lib
COPY --from=libpam-build /usr/lib/security /opt/sysroot/usr/lib/security
COPY --from=libpam-build /usr/include /opt/sysroot/usr/include
COPY --from=libpam-build /usr/share/pam /opt/sysroot/usr/share
COPY --from=libpam-build /usr/sbin /opt/sysroot/usr/sbin

RUN tar -xvf coreutils.tar.gz -C . --strip-components=1 && rm -rf coreutils.tar.gz
RUN cargo install \
        --config /etc/cargo/config.toml \
        --locked \
        --path . \
        --bins \
        --root /build/coreutils

FROM scratch
COPY --from=builder /build/coreutils/bin/* /build/usr/bin
