FROM BUILDER_IMAGE

# build

FROM ubuntu:18.04

RUN apt update --fix-missing && \
    apt-get upgrade -y && \
    apt install -y ca-certificates && \
    apt install -y tzdata && \
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

COPY --from=0 /usr/src/executable-file-name /bin/

ENTRYPOINT ["/bin/executable-file-name", "serve"]
