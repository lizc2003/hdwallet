package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

func DeployContract(opts *bind.TransactOpts, backend bind.ContractBackend, jsonAbi string, bytecode string, params ...interface{}) (common.Address, *types.Transaction, error) {
	parsed, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return common.Address{}, nil, err
	}
	address, tx, _, err := bind.DeployContract(opts, parsed, common.FromHex(bytecode), backend, params...)
	if err != nil {
		return common.Address{}, nil, err
	}
	return address, tx, nil
}

func EstimateContractMethodGas(param TransactBaseParam, backend bind.ContractBackend, contractAddress common.Address, input []byte) (int64, error) {
	err := EnsureTransactGasPrice(backend, &param)
	if err != nil {
		return 0, err
	}

	ethValue := param.EthValue
	if ethValue == nil {
		ethValue = big.NewInt(0)
	}
	msg := ethereum.CallMsg{From: param.From, To: &contractAddress,
		GasPrice:  param.GasPrice,
		GasFeeCap: param.GasFeeCap,
		GasTipCap: param.GasTipCap,
		Value:     ethValue, Data: input}

	gasLimit, err := backend.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	return int64(gasLimit), nil
}
