ARG BUSYBOX_VER="1.38.0"

FROM busybox:$BUSYBOX_VER-musl AS busybox-build
FROM debian:stable-slim AS initramfs-build

WORKDIR /build
RUN bash -c 'mkdir -p /build/{usr/sbin,usr/share/manifest.d,newroot,dev,proc,sys}'

COPY --from=busybox-build /bin/busybox /build/usr/sbin/busybox

# This creates a confusing initramfs environment as agent is renamed to init, however this is a hard requirement of switch_root
COPY --from=initd /app/init /build/init
COPY --from=initd /app/init /build/usr/sbin/init
COPY ./usr/share/manifest.d/init.yaml /build/usr/share/manifest.d/init.yaml

RUN chmod +x usr/sbin/busybox && \
    bash -c 'for prog in $(./usr/sbin/busybox --list); do ln -s ./busybox usr/sbin/$prog; done'

FROM scratch

COPY --from=initramfs-build /build /

SHELL ["/usr/sbin/busybox", "sh"]

ENV PATH="/usr/sbin:/usr/bin"
ENV HOME="/"

CMD ["/usr/sbin/busybox", "sh"]
