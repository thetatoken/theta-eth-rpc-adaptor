FROM golang:1.18

VOLUME ["/go"]

WORKDIR /go/src/github.com/thetatoken/theta-eth-rpc-adaptor/

ENV GOPATH=/go

ENV CGO_ENABLED=1 

ENV GO111MODULE=on

CMD ["/go/src/github.com/thetatoken/theta-eth-rpc-adaptor/integration/docker/build/start.sh"]



