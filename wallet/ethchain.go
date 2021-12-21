package wallet

import (
	"github.com/ethereum/go-ethereum/params"
	"math/big"
)

var (
	BscChainConfig          = &params.ChainConfig{}
	BscTestnetChainConfig   = &params.ChainConfig{}
	MaticChainConfig        = &params.ChainConfig{}
	MaticTestnetChainConfig = &params.ChainConfig{}
)

func init() {
	*MaticChainConfig = *params.AllEthashProtocolChanges
	MaticChainConfig.ChainID = big.NewInt(ChainMatic)
	*MaticTestnetChainConfig = *params.AllEthashProtocolChanges
	MaticTestnetChainConfig.ChainID = big.NewInt(ChainMaticTestnet)
	*BscChainConfig = *params.AllEthashProtocolChanges
	BscChainConfig.ChainID = big.NewInt(ChainBsc)
	*BscTestnetChainConfig = *params.AllEthashProtocolChanges
	BscTestnetChainConfig.ChainID = big.NewInt(ChainBscTestnet)
}
