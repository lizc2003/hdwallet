package eth

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lizc2003/hdwallet/wallet"
	"math/big"
)

type TransactBaseParam struct {
	From      common.Address
	EthValue  *big.Int
	GasPrice  *big.Int
	GasFeeCap *big.Int
	GasTipCap *big.Int
}

func EnsureTransactGasPrice(backend bind.ContractBackend, param *TransactBaseParam) error {
	head, err := backend.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}

	if head.BaseFee == nil {
		if param.GasPrice == nil {
			price, err := backend.SuggestGasPrice(context.Background())
			if err != nil {
				return err
			}
			param.GasPrice = price
		}
	} else {
		if param.GasTipCap == nil {
			tip, err := backend.SuggestGasTipCap(context.Background())
			if err != nil {
				return err
			}
			param.GasTipCap = tip
		}
		if param.GasFeeCap == nil {
			gasFeeCap := new(big.Int).Add(
				param.GasTipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			param.GasFeeCap = gasFeeCap
		}
		if param.GasFeeCap.Cmp(param.GasTipCap) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", param.GasFeeCap, param.GasTipCap)
		}
	}
	return nil
}

func MakeTransactOpts(w *wallet.EthWallet, param TransactBaseParam, gasLimit int64, nonce int64) (*bind.TransactOpts, error) {
	var theNonce *big.Int
	if nonce >= 0 {
		theNonce = big.NewInt(nonce)
	}

	if gasLimit < 0 {
		gasLimit = 0
	}

	txOpts := &bind.TransactOpts{
		From:  param.From,
		Nonce: theNonce,
		Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			return w.SignTx(tx)
		},
		Value:     param.EthValue,
		GasPrice:  param.GasPrice,
		GasFeeCap: param.GasFeeCap,
		GasTipCap: param.GasTipCap,
		GasLimit:  uint64(gasLimit),
		Context:   context.Background(),
	}
	return txOpts, nil
}

func TransferEther(opts *bind.TransactOpts, backend bind.ContractBackend, addressTo common.Address) (*types.Transaction, error) {
	var nonce uint64
	if opts.Nonce != nil {
		nonce = opts.Nonce.Uint64()
	} else {
		tmp, err := backend.PendingNonceAt(context.Background(), opts.From)
		if err != nil {
			return nil, err
		}
		nonce = tmp
	}

	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		gasLimit = wallet.EtherTransferGas
	}

	param := TransactBaseParam{GasPrice: opts.GasPrice,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
	}
	err := EnsureTransactGasPrice(backend, &param)
	if err != nil {
		return nil, err
	}
	opts.GasPrice = param.GasPrice
	opts.GasFeeCap = param.GasFeeCap
	opts.GasTipCap = param.GasTipCap

	var tx *types.Transaction
	var input []byte
	if opts.GasFeeCap == nil {
		baseTx := &types.LegacyTx{
			Nonce:    nonce,
			To:       &addressTo,
			GasPrice: opts.GasPrice,
			Gas:      gasLimit,
			Value:    opts.Value,
			Data:     input,
		}
		tx = types.NewTx(baseTx)
	} else {
		baseTx := &types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &addressTo,
			GasFeeCap: opts.GasFeeCap,
			GasTipCap: opts.GasTipCap,
			Gas:       gasLimit,
			Value:     opts.Value,
			Data:      input,
		}
		tx = types.NewTx(baseTx)
	}

	signedTx, err := opts.Signer(opts.From, tx)
	if err != nil {
		return nil, err
	}

	if opts.NoSend {
		return signedTx, nil
	}

	err = backend.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
