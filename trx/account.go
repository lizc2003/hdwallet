package trx

import (
	"github.com/lizc2003/gotron-sdk/pkg/client"
	"github.com/lizc2003/gotron-sdk/pkg/proto/core"
	"github.com/lizc2003/hdwallet/wallet"
)

func FreezeEnergyBalance(w *wallet.TrxWallet, client *client.GrpcClient, delegateTo string, frozenBalance int64) (string, error) {
	txExt, err := client.FreezeBalance(w.DeriveAddress(), delegateTo, core.ResourceCode_ENERGY, frozenBalance)
	if err != nil {
		return "", err
	}
	return SignAndSendTx(w, client, txExt)
}

func UnfreezeEnergyBalance(w *wallet.TrxWallet, client *client.GrpcClient, delegateTo string) (string, error) {
	txExt, err := client.UnfreezeBalance(w.DeriveAddress(), delegateTo, core.ResourceCode_ENERGY)
	if err != nil {
		return "", err
	}
	return SignAndSendTx(w, client, txExt)
}

func FreezeBandwidthBalance(w *wallet.TrxWallet, client *client.GrpcClient, delegateTo string, frozenBalance int64) (string, error) {
	txExt, err := client.FreezeBalance(w.DeriveAddress(), delegateTo, core.ResourceCode_BANDWIDTH, frozenBalance)
	if err != nil {
		return "", err
	}
	return SignAndSendTx(w, client, txExt)
}

func UnfreezeBandwidthBalance(w *wallet.TrxWallet, client *client.GrpcClient, delegateTo string) (string, error) {
	txExt, err := client.UnfreezeBalance(w.DeriveAddress(), delegateTo, core.ResourceCode_BANDWIDTH)
	if err != nil {
		return "", err
	}
	return SignAndSendTx(w, client, txExt)
}
