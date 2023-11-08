# Cross-compilation issue on Apple Silicon
# https://github.com/mattn/go-sqlite3/issues/992#issuecomment-1002559704

FROM golang:1.18

VOLUME ["/go"]

WORKDIR /go/src/github.com/thetatoken/theta-eth-rpc-adaptor/

ENV GOPATH=/go

RUN dpkg --add-architecture amd64 \
    && apt-get update \
    && apt-get install -y --no-install-recommends gcc-x86-64-linux-gnu g++-x86-64-linux-gnu libc6-dev-amd64-cross

COPY . .

CMD ["/go/src/github.com/thetatoken/theta-eth-rpc-adaptor/integration/docker/build/start.sh"]



