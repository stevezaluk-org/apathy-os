ARG LIBPAM_VER="1.7.2"
ARG LIBCRYPT_VER="4.5.2"

# Base includes our custom build toolchain
FROM ghcr.io/stevezaluk-org/apathy-os/libcrypt:$LIBCRYPT_VER AS libcrypt-build
FROM ghcr.io/stevezaluk-org/apathy-os/base:latest AS builder
ARG LIBPAM_VER
ARG LIBCRYPT_VER

ADD https://github.com/linux-pam/linux-pam/archive/refs/tags/v$LIBPAM_VER.tar.gz libpam.tar.gz

RUN tar -xvf libpam.tar.gz -C . --strip-components=1 && rm -rf libpam.tar.gz

COPY --from=libcrypt-build /usr/lib/* /opt/sysroot/usr/lib
COPY --from=libcrypt-build /usr/include/* /opt/sysroot/usr/include
COPY ./meson-sysroot.ini sysroot.ini

RUN mkdir build && \
    meson setup -Dsecuredir=/usr/lib/security build/ --native-file sysroot.ini

RUN meson compile -C build/ && \
    DESTDIR=/build/libpam/out meson install -C build

FROM scratch
COPY --from=builder /build/libpam/out /
