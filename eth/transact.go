package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/lizc2003/hdwallet/wallet"
	"math/big"
)

type TransactBaseParam struct {
	From      common.Address
	GasPrice  *big.Int
	GasFeeCap *big.Int
	GasTipCap *big.Int
	EthValue  *big.Int
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
