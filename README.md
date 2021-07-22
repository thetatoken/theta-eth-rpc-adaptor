# theta-eth-rpc-adaptor

The `theta-eth-rpc-adaptor` project is aiming to provide an adaptor which translates the Theta RPC interface to the Ethereum RPC interface. Please find the currently supported Ethereum RPC APIs [here](https://github.com/thetatoken/theta-eth-rpc-adaptor#rpc-apis).

## Setup

### On macOS

Install Go and set environment variables `GOPATH` , `GOBIN`, and `PATH`. The current code base should compile with **Go 1.14.2**. On macOS, install Go with the following command

```
brew install go@1.14.2
brew link go@1.14.2 --force
```

First clone the `theta` repo following the steps below. Then, clone this repo into your `$GOPATH`. The path should look like this: `$GOPATH/src/github.com/thetatoken/theta`

```
mkdir -p $GOPATH/src/github.com/thetatoken 
cd $GOPATH/src/github.com/thetatoken
git clone https://github.com/thetatoken/theta-protocol-ledger.git $GOPATH/src/github.com/thetatoken/theta
cd theta
git checkout theta3.0-rpc-compatibility
export GO111MODULE=on
make install
```

Next, clone the `theta-eth-rpc-adaptor` repo:

```
cd $GOPATH/src/github.com/thetatoken
git clone https://github.com/thetatoken/theta-eth-rpc-adaptor
```

## Build and Install

### Build the binary under macOS or Linux
Following the steps below to build the `theta-eth-rpc-adaptor` binary and copy it into your `$GOPATH/bin`.

```
export THETA_ETH_RPC_ADAPTOR_HOME=$GOPATH/src/github.com/thetatoken/theta-eth-rpc-adaptor
cd $THETA_ETH_RPC_ADAPTOR_HOME
export GO111MODULE=on
make install
```

### Cross compilation for Windows
On a macOS machine, the following command should build the `theta-eth-rpc-adaptor.exe` binary under `build/windows/`

```
make windows
```

## Run the Adaptor with a local Theta private testnet

First, run a private testnet Theta node with its RPC port opened at 16888:

```
cd $THETA_HOME
cp -r ./integration/privatenet ../privatenet
mkdir ~/.thetacli
cp -r ./integration/privatenet/thetacli/* ~/.thetacli/
chmod 700 ~/.thetacli/keys/encrypted

theta start --config=../privatenet/node --password=qwertyuiop
```

Then, open another terminal, create the config folder for the RPC adaptor

```
export THETA_ETH_RPC_ADAPTOR_HOME=$GOPATH/src/github.com/thetatoken/theta-eth-rpc-adaptor
cd $THETA_ETH_RPC_ADAPTOR_HOME
mkdir ../privatenet/eth-rpc-adaptor
```

Use you favorite editor to open file `../privatenet/eth-rpc-adaptor/config.yaml`, paste in the follow content, save and close the file:

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

Then, launch the adaptor binary with the following command:

```
cd $THETA_ETH_RPC_ADAPTOR_HOME
mkdir ../privatenet/eth-rpc-adaptor
theta-eth-rpc-adaptor start --config=../privatenet/eth-rpc-adaptor
```

## RPC APIs

The RPC APIs should conform to the Ethereum JSON RPC API standard: https://eth.wiki/json-rpc/API. We currently support the following Ethereum RPC APIs:

```
eth_chainId
eth_syncing
eth_accounts
eth_protocolVersion
eth_getBlockByHash
eth_getBlockByNumber
eth_blockNumber
eth_getUncleByBlockHashAndIndex
eth_getTransactionByHash
eth_getTransactionByBlockNumberAndIndex
eth_getTransactionByBlockHashAndIndex
eth_getBlockTransactionCountByHash
eth_getTransactionReceipt
eth_getBalance
eth_getStorageAt
eth_getCode
eth_getTransactionCount
eth_getLogs
eth_getBlockTransactionCountByNumber
eth_call
eth_gasPrice
eth_estimateGas
eth_sendRawTransaction
eth_sendTransaction
net_version
```

The following examples demonstrate how to interact with the RPC APIs using the `curl` command:

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
