package ethrpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/thetatoken/theta-eth-rpc-adaptor/common"
)

var testTypeData = `{"types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"version","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"Person":[{"name":"name","type":"string"},{"name":"test","type":"uint8"},{"name":"wallet","type":"address"}],"Mail":[{"name":"from","type":"Person"},{"name":"to","type":"Person"},{"name":"contents","type":"string"}]},"primaryType":"Mail","domain":{"name":"Ether Mail","version":"1","chainId":"1","verifyingContract":"0xCCCcccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC"},"message":{"from":{"name":"Cow","test":"3","wallet":"0xcD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826"},"to":{"name":"Bob","wallet":"0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB","test":"2"},"contents":"Hello, Bob!"}}`

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
	logger.Infof("jlog1 typedData is %+v, chain ID is %v \n", typedData, typedData.Domain.ChainId)

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return
	}
	logger.Infof("jlog2 domainSeparator: %+v \n", domainSeparator)

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return
	}
	logger.Infof("jlog3 typedDataHash: %+v \n", typedDataHash)

	dataStr := fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	rawData := []byte(dataStr)
	logger.Infof("jlog4 dataStr: %s \n", dataStr)
	result, err = e.Sign(ctx, strings.ToLower(address), rawData)
	if err != nil {
		return
	}
	// result = hex.EncodeToString(signature.ToBytes())
	logger.Infof("jlog6 result: %s \n", result)
	return result, nil
}

/*
// SignTypedData signs EIP-712 conformant typed data
// hash = keccak256("\x19${byteVersion}${domainSeparator}${hashStruct(message)}")
// It returns
// - the signature,
// - and/or any error
func (api *SignerAPI) SignTypedData(ctx context.Context, addr common.MixedcaseAddress, typedData TypedData) (hexutil.Bytes, error) {
	signature, _, err := api.signTypedData(ctx, addr, typedData, nil)
	return signature, err
}

// signTypedData is identical to the capitalized version, except that it also returns the hash (preimage)
// - the signature preimage (hash)
func (api *SignerAPI) signTypedData(ctx context.Context, addr common.MixedcaseAddress,
	typedData TypedData, validationMessages *apitypes.ValidationMessages) (hexutil.Bytes, hexutil.Bytes, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, nil, err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, nil, err
	}
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(rawData)
	messages, err := typedData.Format()
	if err != nil {
		return nil, nil, err
	}
	req := &SignDataRequest{
		ContentType: DataTyped.Mime,
		Rawdata:     rawData,
		Messages:    messages,
		Hash:        sighash,
		Address:     addr}
	if validationMessages != nil {
		req.Callinfo = validationMessages.Messages
	}
	signature, err := api.sign(req, true)
	if err != nil {
		api.UI.ShowError(err.Error())
		return nil, nil, err
	}
	return signature, sighash, nil
}

// HashStruct generates a keccak256 hash of the encoding of the provided data
func (typedData *TypedData) HashStruct(primaryType string, data TypedDataMessage) (hexutil.Bytes, error) {
	encodedData, err := typedData.EncodeData(primaryType, data, 1)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encodedData), nil
}
*/
