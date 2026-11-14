ARG UTILLINUX_VER="2.42.2"
ARG LIBPAM_VER="1.7.2"
ARG LIBCRYPT_VER="4.5.2"

FROM ghcr.io/stevezaluk-org/apathy-os/libcrypt:$LIBCRYPT_VER AS libcrypt-build
FROM ghcr.io/stevezaluk-org/apathy-os/libpam:$LIBPAM_VER AS libpam-build
FROM ghcr.io/stevezaluk-org/apathy-os/base:latest AS builder
ARG UTILLINUX_VER
ARG LIBPAM_VER
ARG LIBCRYPT_VER

WORKDIR /build/util-linux

ADD https://github.com/util-linux/util-linux/archive/refs/tags/v$UTILLINUX_VER.tar.gz util-linux.tar.gz

RUN tar -xvf util-linux.tar.gz -C . --strip-components=1 && rm -rf util-linux.tar.gz

RUN apt-get update && \
    apt-get install bison flex libtool make automake meson gettext autopoint pkg-config -y

COPY --from=libcrypt-build /usr/lib/* /opt/sysroot/usr/lib
COPY --from=libcrypt-build /usr/include /opt/sysroot/usr/include
COPY --from=libpam-build /usr/lib/x86_64-linux-gnu/ /opt/sysroot/usr/lib
COPY --from=libpam-build /usr/lib/security /opt/sysroot/usr/lib/security
COPY --from=libpam-build /usr/include /opt/sysroot/usr/include
COPY --from=libpam-build /usr/share/pam /opt/sysroot/usr/share
COPY --from=libpam-build /usr/sbin /opt/sysroot/usr/sbin

# Python disabled here so that we don't try and build python bindings for libmount
# libaudit doesn't officially support cross-compiling (lame). See: https://github.com/linux-audit/audit-userspace/tree/master#cross-compiling
RUN meson setup build \
    --prefix=/usr \
    -Dbuild-python=disabled \
    -Daudit=disabled \
    --native-file /opt/conf/meson-sysroot.ini

RUN DESTDIR=/build/util-linux/out meson install -C build/
