FROM golang:1.12

ARG CI_JOB_TOKEN
ADD . /srv/codes

RUN git config --global credential.helper store && \
    echo "https://gitlab-ci-token:${CI_JOB_TOKEN}@git.cafebazaar.ir" >> ~/.git-credentials && \
    export GOPATH=/srv/codes/__gopath__ && \
    export PATH=${PATH%:/home/go/bin} && \
    export PATH=$PATH:/$GOPATH/bin && \
    if [ -d $GOPATH/src/git.cafebazaar.ir/arcana261/golang-boilerplate ]; then rm -rf $GOPATH/src/git.cafebazaar.ir/arcana261/golang-boilerplate; fi && \
    mkdir -p $GOPATH/src/git.cafebazaar.ir/arcana261/golang-boilerplate && \
    for file in $(find /srv/codes -maxdepth 1 ! -name '__gopath__' ! -name '.'); do cp -rf $file $GOPATH/src/git.cafebazaar.ir/arcana261/golang-boilerplate; done && \
    cd $GOPATH/src/git.cafebazaar.ir/arcana261/golang-boilerplate && \
    make dependencies && \
    make postviewd

# build

FROM ubuntu:18.04

RUN apt update --fix-missing && \
    apt-get upgrade -y && \
    apt install -y ca-certificates && \
    apt install -y tzdata && \
    ln -sf /usr/share/zoneinfo/UTC /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata && \
    apt-get clean

COPY --from=0 /srv/codes/__gopath__/git.cafebazaar.ir/arcana261/golang-boilerplate/postviewd /bin/

ENTRYPOINT ["/bin/postviewd", "serve"]
