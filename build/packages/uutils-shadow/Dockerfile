ARG LIBPAM_VER="1.7.2"
ARG LIBCRYPT_VER="4.5.2"
ARG SHADOW_VER="0.2.1"

FROM ghcr.io/stevezaluk-org/apathy-os/libcrypt:$LIBCRYPT_VER AS libcrypt-build
FROM ghcr.io/stevezaluk-org/apathy-os/libpam:$LIBPAM_VER AS libpam-build
FROM ghcr.io/stevezaluk-org/apathy-os/base:latest AS builder
ARG LIBPAM_VER
ARG LIBCRYPT_VER
ARG SHADOW_VER

WORKDIR /build/shadow-rs

ADD https://github.com/uutils/shadow/archive/refs/tags/v$SHADOW_VER.tar.gz shadow-rs.tar.gz

COPY --from=libcrypt-build /usr/lib/* /opt/sysroot/usr/lib
COPY --from=libcrypt-build /usr/include /opt/sysroot/usr/include
COPY --from=libpam-build /usr/lib/x86_64-linux-gnu/ /opt/sysroot/usr/lib
COPY --from=libpam-build /usr/lib/security /opt/sysroot/usr/lib/security
COPY --from=libpam-build /usr/include /opt/sysroot/usr/include
COPY --from=libpam-build /usr/share/pam /opt/sysroot/usr/share
COPY --from=libpam-build /usr/sbin /opt/sysroot/usr/sbin

RUN tar -xvf shadow-rs.tar.gz -C . --strip-components=1 && rm -rf shadow-rs.tar.gz
RUN make install PREFIX=/build/shadow-rs/bin

FROM scratch
COPY --from=builder /build/shadow-rs/bin/* /usr/bin
