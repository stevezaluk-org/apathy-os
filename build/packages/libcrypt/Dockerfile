ARG LIBCRYPT_VER="4.5.2"

FROM ghcr.io/stevezaluk-org/apathy-os/base:latest AS builder
ARG LIBCRYPT_VER

ADD https://github.com/besser82/libxcrypt/archive/refs/tags/v$LIBCRYPT_VER.tar.gz libcrypt.tar.gz

RUN apt-get update -y && \
    apt-get install tar autoconf libtool make -y

RUN tar -xvf libcrypt.tar.gz -C . --strip-components=1 && rm -rf libcrypt.tar.gz

RUN ./autogen.sh
RUN ./configure --with-sysroot=/opt/sysroot \
                --prefix=/usr

RUN make -j$(nproc) && \
    DESTDIR="/build/libcrypt/out" make install

FROM scratch

COPY --from=builder /build/libcrypt/out /
