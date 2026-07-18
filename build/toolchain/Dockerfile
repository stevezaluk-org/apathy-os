FROM debian:stable-slim AS crosstools-build
ARG CROSSTOOLS_VER="1.28.0"

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

COPY ./config .config

RUN ct-ng build