package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
