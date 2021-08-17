package ethrpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta/common/hexutil"
	// "golang.org/x/crypto/sha3"
)

const (
	MimetypeTypedData  = "data/typed"
	MimetypeTextPlain  = "text/plain"
	personalSignPrefix = "\x19Ethereum Signed Message:\n"
)

// ------------------------------- eth_sign -----------------------------------
func (e *EthRPCService) Sign(ctx context.Context, account string, message string) (result string, err error) {
	logger.Infof("eth_sign called, account: %s, message: %v \n", account, message)
	msgBytes, _ := hexutil.Decode(message)
	signhash := []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(msgBytes), msgBytes))
	signature, err := common.SignRawBytes(strings.ToLower(account), signhash)
	if err != nil {
		return
	}
	result = hexutil.Encode(signature.ToBytes())
	return result, nil
}
