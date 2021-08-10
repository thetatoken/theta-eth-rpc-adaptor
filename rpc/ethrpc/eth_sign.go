package ethrpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
	"github.com/thetatoken/theta/common/hexutil"
	"golang.org/x/crypto/sha3"
)

const (
	MimetypeTypedData  = "data/typed"
	MimetypeTextPlain  = "text/plain"
	personalSignPrefix = "\x19Ethereum Signed Message:\n"
)

// ------------------------------- eth_sign -----------------------------------
func (e *EthRPCService) Sign(ctx context.Context, account string, message hexutil.Bytes) (result string, err error) {
	logger.Infof("eth_sign called, account: %s, message: %s \n", account, message)
	// sighash, _ := TextAndHash(message)
	sighash := message
	signature, err := common.SignRawBytes(strings.ToLower(account), sighash)
	if err != nil {
		return
	}
	if signature.ToBytes()[64] < 2 {
		signature.ToBytes()[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	}
	result = hexutil.Encode(signature.ToBytes())
	logger.Infof("jlog7 result : %s \n", result)
	return result, nil
}

func TextAndHash(data []byte) ([]byte, string) {
	msg := fmt.Sprintf("%s%d%s", personalSignPrefix, len(data), string(data))
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(msg))
	return hasher.Sum(nil), msg
}
