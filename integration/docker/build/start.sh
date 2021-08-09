#!/bin/bash

echo "Building binaries..."

set -e
set -x

GOBIN=/usr/local/go/bin/go

$GOBIN build -o ./build/linux/theta-eth-rpc-adaptor ./cmd/theta-eth-rpc-adaptor

set +x 

echo "Done."



