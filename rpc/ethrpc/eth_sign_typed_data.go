package ethrpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/thetatoken/theta/common/hexutil"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
)

// ------------------------------- eth_signTypedData -----------------------------------
func (e *EthRPCService) SignTypedData(ctx context.Context, address string, typedDataObj common.TypedDataPara) (result string, err error) {
	logger.Infof("eth_signTypedData called, address: %s, typedData: %+v \n", address, typedDataObj)

	typedData := common.TypedData{
		Types:       typedDataObj.Types,
		PrimaryType: typedDataObj.PrimaryType,
		Domain: common.TypedDataDomain{
			Name:              typedDataObj.Domain.Name,
			Version:           typedDataObj.Domain.Version,
			VerifyingContract: typedDataObj.Domain.VerifyingContract,
			Salt:              typedDataObj.Domain.Salt,
			ChainId:           common.NewHexOrDecimal256(typedDataObj.Domain.ChainId),
		},
		Message: typedDataObj.Message,
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return
	}
	dataStr := fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	rawData := []byte(dataStr)
	signature, err := common.SignRawBytes(strings.ToLower(address), rawData)
	if err != nil {
		return
	}
	signature.ToBytes()[64] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	result = hexutil.Encode(signature.ToBytes())
	return result, nil
}
