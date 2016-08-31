# Dockerfile used for cross compiling Windows and Linux binaries from Linux.
#
# This is needed because os/user uses CGO, thus cross compiling needs
# CGO_ENABLED=1. Otherwise the following error will happen:
#
#   'user: Current not implemented on linux/amd64'
#

FROM golang
MAINTAINER Calin Seciu

RUN apt-get update && apt-get install -y --no-install-recommends \
        pkg-config \
        cmake \
        mingw-w64 \
    && rm -rf /var/lib/apt/lists/*

# Install Go.
#   1) 1.4 for bootstrap.
ENV GOROOT_BOOTSTRAP /go1.4
RUN (curl -sSL https://golang.org/dl/go1.4.linux-amd64.tar.gz | tar -vxz -C /tmp) && \
	mv /tmp/go $GOROOT_BOOTSTRAP

# Compile the Windows 64-bit toolchain
RUN cd /usr/local/go/src/ && \
    env CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC_FOR_TARGET="x86_64-w64-mingw32-gcc" ./make.bash --no-clean

COPY . "$GOPATH/src/github.com/clns/izip-client"

WORKDIR "$GOPATH/src/github.com/clns/izip-client"

RUN go get -d ./... && (GOOS=windows go get -d ./... || true)
