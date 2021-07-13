package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"github.com/lizc2003/hdwallet/wallet"
)

type TransactBaseParam struct {
	From     common.Address
	GasPrice *big.Int
	EthValue *big.Int
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
		Value:    param.EthValue,
		GasPrice: param.GasPrice,
		GasLimit: uint64(gasLimit),
		Context: context.Background(),
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

	tx := types.NewTransaction(nonce, addressTo, opts.Value, gasLimit, opts.GasPrice, nil)
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
