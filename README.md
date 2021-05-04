# theta-eth-rpc-adaptor

## Setup

### On macOS

Install Go and set environment variables `GOPATH` , `GOBIN`, and `PATH`. The current code base should compile with **Go 1.12.1**. On macOS, install Go with the following command

```
brew install go@1.12.1
brew link go@1.12.1 --force
```

First clone the `theta` repo following the steps below. Then, clone this repo into your `$GOPATH`. The path should look like this: `$GOPATH/src/github.com/thetatoken/edgecore`

```
git clone https://github.com/thetatoken/theta-protocol-ledger.git $GOPATH/src/github.com/thetatoken/theta
export THETA_HOME=$GOPATH/src/github.com/thetatoken/theta
cd $THETA_HOME

git clone https://github.com/thetatoken/theta-eth-rpc-adaptor
export THETA_ETH_RPC_ADAPTOR_HOME=$GOPATH/src/github.com/thetatoken/theta-eth-rpc-adaptor
cd $THETA_ETH_RPC_ADAPTOR_HOME
```

## Build and Install

### Build the binary under macOS or Linux
This should build the `theta-eth-rpc-adaptor` binary and copy it into your `$GOPATH/bin`.

```
export GO111MODULE=on
make install
```

### Cross compilation for Windows
On a macOS machine, the following command should build the `theta-eth-rpc-adaptor.exe` binary under `build/windows/`

```
make windows
```

## Run the Edge Core
Launch the edge core binary with the following command

```
theta-eth-rpc-adaptor start --config=<CONFIG_FOLDER>
```

## RPC APIs

The RPC APIs should conform to the Ethereum JSON RPC API standard: https://eth.wiki/json-rpc/API

```
# Query version
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[],"id":67}' http://localhost:18888/rpc

# Query synchronization status
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' http://localhost:18888/rpc

# Query block number
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://localhost:18888/rpc

# Query account balance (should return an integer which represents the current TFuel balance in wei)
curl -X POST --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"],"id":1}' http://localhost:17888/rpc
```
