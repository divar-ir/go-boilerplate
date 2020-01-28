FROM build_image_tag as build_image

FROM ubuntu:18.04

RUN apt update --fix-missing && \
    apt-get upgrade -y && \
    apt install -y ca-certificates && \
    apt install -y tzdata && \
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

COPY --from=build_image /srv/build/appdetaild /bin/

ENTRYPOINT ["/bin/appdetaild", "serve"]
