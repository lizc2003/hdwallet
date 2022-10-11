package trx

import (
	"github.com/lizc2003/gotron-sdk/pkg/client"
	"github.com/lizc2003/hdwallet/wallet"
)

func TransferTrx(w *wallet.TrxWallet, client *client.GrpcClient, toAddr string, amount int64) (string, error) {
	txExt, err := client.Transfer(w.DeriveAddress(), toAddr, amount)
	if err != nil {
		return "", err
	}

	return SignAndSendTx(w, client, txExt)
}
