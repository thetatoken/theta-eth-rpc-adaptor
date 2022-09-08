module github.com/thetatoken/theta-eth-rpc-adaptor

require (
	github.com/dgraph-io/badger v1.6.1 // indirect
	github.com/ethereum/go-ethereum v1.9.23
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pborman/uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.13.0
	//github.com/stretchr/testify v1.7.0
	github.com/syndtr/goleveldb v1.0.1-0.20210305035536-64b5b1c73954 // indirect
	github.com/thetatoken/theta v0.0.0
	github.com/thetatoken/theta/common v0.0.0
	github.com/ybbus/jsonrpc v1.1.1
)

replace github.com/thetatoken/theta v0.0.0 => ../theta

replace github.com/thetatoken/theta/common v0.0.0 => ../theta/common

replace github.com/thetatoken/theta/rpc/lib/rpc-codec/jsonrpc2 v0.0.0 => ../theta/rpc/lib/rpc-codec/jsonrpc2/

replace github.com/ethereum/go-ethereum => github.com/ethereum/go-ethereum v1.9.9

go 1.13
