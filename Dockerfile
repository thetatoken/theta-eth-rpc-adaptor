FROM golang:1.14.2 as build
# RUN apt-get update && apt-get install -y curl make gcc g++ git
# ENV CGO_ENABLED=1 
ENV GO111MODULE=on
# ENV GIT_TERMINAL_PROMPT=1
ENV THETA_TOKEN_HOME=$GOPATH/src/github.com/thetatoken
WORKDIR $THETA_TOKEN_HOME/theta
RUN git clone https://github.com/thetatoken/theta-protocol-ledger.git .
RUN git checkout testnet
RUN make install
WORKDIR $THETA_TOKEN_HOME/theta-eth-rpc-adaptor
COPY ./ .
RUN make install
# FROM alpine:latest
# RUN apk add --no-cache ca-certificates
ENV PATH=$GOPATH/bin:/usr/local/go/bin:/usr/local/bin:$PATH

RUN apt-get update \
  && apt-get install -y vim \
  && DEBIAN_FRONTEND=noninteractive apt-get install -y \
  net-tools \
  && apt-get clean \
  && rm -rf /var/lib/apt/lists/*

CMD [ "theta-eth-rpc-adaptor", "start", "--config=." ]
EXPOSE 18888