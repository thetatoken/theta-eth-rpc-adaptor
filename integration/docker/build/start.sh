#!/bin/bash

echo "Building binaries..."

set -e
set -x

GOBIN=/usr/local/go/bin/go

CGO_ENABLED=1 GOARCH=amd64 CC=x86_64-linux-gnu-gcc CXX=g++-x86-64-linux-gnu $GOBIN build -o ./build/linux/theta-eth-rpc-adaptor ./cmd/theta-eth-rpc-adaptor

set +x 

echo "Done."



