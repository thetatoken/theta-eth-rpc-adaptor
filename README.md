# theta-eth-rpc-adaptor

The `theta-eth-rpc-adaptor` project is aiming to provide an adaptor which translates the Theta RPC interface to the Ethereum RPC interface.

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

## Run the Adaptor

To run the adaptor, you'd first need to run a Theta node on the same machine with its RPC port opened at 16888. Then, in another terminal, launch the adaptor binary with the following command, assuming the `config.yaml` file is placed under the `<CONFIG_FOLDER>/` folder:

```
theta-eth-rpc-adaptor start --config=<CONFIG_FOLDER>
```

Below is an example `config.yaml` file

```
theta:
  rpcEndpoint: "http://127.0.0.1:16888/rpc"
rpc:
  enabled: true
  httpAddress: "127.0.0.1"
  httpPort: 18888
  wsAddress: "127.0.0.1"
  wsPort: 18889
  timeoutSecs: 600 
  maxConnections: 2048
log:
  levels: "*:debug"
```

For example, you can change the above `theta.rpcEnpoint` to a remote Theta RPC endpoint, or change `rpc.httpAddress` to "0.0.0.0" so the adaptor is accessor from remote IP addresses.

## RPC APIs

The RPC APIs should conform to the Ethereum JSON RPC API standard: https://eth.wiki/json-rpc/API

```
# Query version
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_protocolVersion","params":[],"id":67}' http://localhost:18888/rpc

# Query synchronization status
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_syncing","params":[],"id":1}' http://localhost:18888/rpc

# Query block number
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' http://localhost:18888/rpc

# Query account TFuel balance (should return an integer which represents the current TFuel balance in wei)
curl -X POST -H 'Content-Type: application/json' --data '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x407d73d8a49eeb85d32cf465507dd71d507100c1", "latest"],"id":1}' http://localhost:18888/rpc
```
