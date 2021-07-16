package eth

import (
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
)

var BigIntEthGWei = big.NewInt(1e9)

func HexToAddress(addr string) (common.Address, error) {
	if !common.IsHexAddress(addr) {
		return common.Address{}, errors.New("invalid address")
	}

	return common.HexToAddress(addr), nil
}

func WeiToGwei(v *big.Int) int64 {
	return big.NewInt(0).Div(v, BigIntEthGWei).Int64()
}

func GweiToWei(v int64) *big.Int {
	return big.NewInt(0).Mul(big.NewInt(v), BigIntEthGWei)
}

func CalcEthFee(gasPrice *big.Int, gas int64) int64 {
	return WeiToGwei(big.NewInt(0).Mul(big.NewInt(gas), gasPrice))
}

func SerializeTransaction(tx *types.Transaction) (string, error) {
	data, err := tx.MarshalBinary()
	if err != nil {
		return "", err
	}
	return hexutil.Encode(data), nil
}
