package trx

import "C"
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/lizc2003/gotron-sdk/pkg/client"
	"github.com/lizc2003/gotron-sdk/pkg/common"
	"github.com/lizc2003/gotron-sdk/pkg/proto/api"
	"github.com/lizc2003/gotron-sdk/pkg/proto/core"
	"github.com/lizc2003/hdwallet/wallet"
	"google.golang.org/protobuf/proto"
)

type TrxTransaction struct {
	tx   *core.Transaction
	txId string
}

func NewTransaction(txExt *api.TransactionExtention) (*TrxTransaction, error) {
	return &TrxTransaction{tx: txExt.Transaction, txId: hex.EncodeToString(txExt.Txid)}, nil
}

func (this *TrxTransaction) TxHash() ([]byte, error) {
	rawData, err := proto.Marshal(this.tx.GetRawData())
	if err != nil {
		return nil, err
	}
	txHash := sha256.Sum256(rawData)
	return txHash[:], nil
}

func (this *TrxTransaction) Sign(w *wallet.TrxWallet) error {
	txHash, err := this.TxHash()
	if err != nil {
		return err
	}
	signature, err := crypto.Sign(txHash, w.DeriveNativePrivateKey())
	if err != nil {
		return err
	}
	this.tx.Signature = append(this.tx.Signature, signature)
	return nil
}

func (this *TrxTransaction) Send(client *client.GrpcClient) (string, error) {
	result, err := client.Broadcast(this.tx)
	if err != nil {
		return "", err
	}
	if result.Code != api.Return_SUCCESS {
		return "", fmt.Errorf("send transaction fail: %s, %s", result.Code.String(), string(result.GetMessage()))
	}

	h, _ := this.TxHash()
	return common.BytesToHexString(h), nil
}

func SignAndSendTx(w *wallet.TrxWallet, client *client.GrpcClient, txExt *api.TransactionExtention) (string, error) {
	tx, err := NewTransaction(txExt)
	if err != nil {
		return "", err
	}
	err = tx.Sign(w)
	if err != nil {
		return "", err
	}
	return tx.Send(client)
}
