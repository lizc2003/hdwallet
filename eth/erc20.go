package eth

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"strings"
)

const ERC20InterfaceABI = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"address","name":"spender","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"from","type":"address"},{"indexed":true,"internalType":"address","name":"to","type":"address"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"address","name":"spender","type":"address"}],"name":"allowance","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"spender","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"approve","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"account","type":"address"}],"name":"balanceOf","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"decimals","outputs":[{"internalType":"uint8","name":"","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"name","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"symbol","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"totalSupply","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transfer","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"sender","type":"address"},{"internalType":"address","name":"recipient","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"transferFrom","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"nonpayable","type":"function"}]`

// Erc20Contract tool for contract abi
type Erc20Contract struct {
	abi             abi.ABI
	contractAddress common.Address
	backend         bind.ContractBackend
	contract        *bind.BoundContract
	opts            *bind.CallOpts
}

func NewErc20Contract(address common.Address, backend bind.ContractBackend) *Erc20Contract {
	parsed, _ := abi.JSON(strings.NewReader(ERC20InterfaceABI))
	c := bind.NewBoundContract(address, parsed, backend, backend, backend)
	return &Erc20Contract{abi: parsed, contractAddress: address, backend: backend, contract: c, opts: &bind.CallOpts{}}
}

func (this *Erc20Contract) TotalSupply() (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := this.contract.Call(this.opts, &out, "totalSupply")
	if err != nil {
		return nil, err
	}
	return *ret0, err
}

func (this *Erc20Contract) Name() (string, error) {
	var (
		ret0 = new(string)
	)
	out := []interface{}{ret0}
	err := this.contract.Call(this.opts, &out, "name")
	if err != nil {
		return "", err
	}
	return *ret0, err
}

func (this *Erc20Contract) Symbol() (string, error) {
	var (
		ret0 = new(string)
	)
	out := []interface{}{ret0}
	err := this.contract.Call(this.opts, &out, "symbol")
	if err != nil {
		return "", err
	}
	return *ret0, err
}

func (this *Erc20Contract) Decimals() (int, error) {
	var (
		ret0 = new(uint8)
	)
	out := []interface{}{ret0}
	err := this.contract.Call(this.opts, &out, "decimals")
	if err != nil {
		return 0, err
	}
	return int(*ret0), err
}

func (this *Erc20Contract) BalanceOf(tokenOwner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := this.contract.Call(this.opts, &out, "balanceOf", tokenOwner)
	if err != nil {
		return nil, err
	}
	return *ret0, err
}

func (this *Erc20Contract) Allowance(tokenOwner common.Address, spender common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := []interface{}{ret0}
	err := this.contract.Call(this.opts, &out, "allowance", tokenOwner, spender)
	if err != nil {
		return nil, err
	}
	return *ret0, err
}

func (this *Erc20Contract) Transfer(opts *bind.TransactOpts, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return this.contract.Transact(opts, "transfer", to, tokens)
}

func (this *Erc20Contract) Approve(opts *bind.TransactOpts, spender common.Address, tokens *big.Int) (*types.Transaction, error) {
	return this.contract.Transact(opts, "approve", spender, tokens)
}

func (this *Erc20Contract) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokens *big.Int) (*types.Transaction, error) {
	return this.contract.Transact(opts, "transferFrom", from, to, tokens)
}

func (this *Erc20Contract) EstimateTransferGas(param TransactBaseParam, to common.Address, tokens *big.Int) (int64, error) {
	input, err := this.abi.Pack("transfer", to, tokens)
	if err != nil {
		return 0, err
	}
	return EstimateContractMethodGas(param, this.backend, this.contractAddress, input)
}

func (this *Erc20Contract) EstimateApproveGas(param TransactBaseParam, spender common.Address, tokens *big.Int) (int64, error) {
	input, err := this.abi.Pack("approve", spender, tokens)
	if err != nil {
		return 0, err
	}
	return EstimateContractMethodGas(param, this.backend, this.contractAddress, input)
}

func (this *Erc20Contract) EstimateTransferFromGas(param TransactBaseParam, from common.Address, to common.Address, tokens *big.Int) (int64, error) {
	input, err := this.abi.Pack("transferFrom", from, to, tokens)
	if err != nil {
		return 0, err
	}
	return EstimateContractMethodGas(param, this.backend, this.contractAddress, input)
}
