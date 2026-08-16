ARG TARGET="x86_64-unknown-linux-gnu"
ARG CROSSTOOLS_VER="1.28.0"
ARG LIBCRYPT_VER="4.5.2"
ARG LIBPAM_VER="1.7.2"

FROM debian:stable-slim AS crosstools-build
ARG TARGET
ARG CROSSTOOLS_VER

RUN apt-get update && \
    apt-get install make bzip2 wget rsync gcc python3 python3-dev python-is-python3 flex bison texinfo xz-utils unzip help2man file patch gawk libtool libtool-bin build-essential libncurses-dev adduser -y

WORKDIR /build/crosstools

ADD http://crosstool-ng.org/download/crosstool-ng/crosstool-ng-1.28.0.tar.bz2 crosstools-ng.tar.bz2

RUN tar -xvf crosstools-ng.tar.bz2 -C . --strip-components=1 && rm -rf crosstools-ng.tar.bz2

RUN CXXFLAGS="-std=c++17" ./configure --prefix=/usr

RUN make -j $(nproc) && \
    make install

WORKDIR /build/toolchain

RUN groupadd --system nonroot && \
    useradd -m -s /bin/bash nonroot -g nonroot

RUN chown -R nonroot:nonroot /build/toolchain

USER nonroot:nonroot

RUN ct-ng $TARGET && \
    ct-ng build

FROM debian:stable-slim AS libcrypt-build
ARG TARGET
ARG LIBCRYPT_VER

RUN apt-get update && \
    apt-get install autoconf libtool make -y

WORKDIR /build/libcrypt

ADD https://github.com/besser82/libxcrypt/archive/refs/tags/v$LIBCRYPT_VER.tar.gz libcrypt.tar.gz

RUN tar -xvf libcrypt.tar.gz -C . --strip-components=1 && rm -rf libcrypt.tar.gz

# We need to ensure that we are actually compiling against our toolchain, so we need to copy this in
COPY --from=crosstools-build /home/nonroot/x-tools/$TARGET/$TARGET/sysroot/ /build/sysroot

ENV CPPFLAGS="-I/build/sysroot/usr/include"
ENV LDFLAGS="-L/build/sysroot/usr/lib -Wl,-rpath-link,/build/sysroot/usr/lib"
ENV PKG_CONFIG_SYSROOT_DIR="/build/sysroot"
ENV PKG_CONFIG_LIBDIR="/build/sysroot/usr/lib/pkgconfig"

RUN ./autogen.sh
RUN ./configure CC="gcc --sysroot=/build/sysroot" \
                CXX="g++ --sysroot=/build/sysroot" \
                --with-sysroot=/build/sysroot \
                --prefix=/usr

RUN make -j$(nproc) && \
    DESTDIR="/build/libcrypt/out" make install

FROM debian:stable-slim AS libpam-build
ARG TARGET
ARG LIBPAM_VER

RUN apt-get update && \
    apt-get install meson pkg-config -y

WORKDIR /build/libpam

ADD https://github.com/linux-pam/linux-pam/archive/refs/tags/v$LIBPAM_VER.tar.gz libpam.tar.gz

RUN tar -xvf libpam.tar.gz -C . --strip-components=1 && rm -rf libpam.tar.gz

# We need to ensure that we are actually compiling against our toolchain, so we need to copy this in
COPY --from=crosstools-build /home/nonroot/x-tools/$TARGET/$TARGET/sysroot/ /build/sysroot
COPY --from=libcrypt-build /build/libcrypt/out/usr/lib/* /build/sysroot/usr/lib
COPY --from=libcrypt-build /build/libcrypt/out/usr/include/* /build/sysroot/usr/include
COPY ./meson-sysroot.ini sysroot.ini

RUN mkdir build && \
    meson setup build/ --native-file sysroot.ini

RUN meson compile -C build/ && \
    DESTDIR=/build/libpam/out meson install -C build

FROM scratch
ARG TARGET

COPY --from=crosstools-build /home/nonroot/x-tools/$TARGET/$TARGET/sysroot/ /
COPY --from=libcrypt-build /build/libcrypt/out/usr/lib/* /usr/lib

COPY --from=libpam-build /build/libpam/out/usr/lib/x86_64-linux-gnu/* /usr/lib
COPY --from=libpam-build /build/libpam/out/usr/share/pam /usr/share
COPY --from=libpam-build /build/libpam/out/usr/sbin /usr/sbin