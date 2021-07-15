package btc

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/lizc2003/hdwallet/wallet"
)

type BtcUnspent struct {
	TxID         string  `json:"txid"`
	Vout         uint32  `json:"vout"`
	ScriptPubKey string  `json:"scriptPubKey"`
	RedeemScript string  `json:"redeemScript,omitempty"`
	Amount       float64 `json:"amount"`
}

type BtcOutput struct {
	Address btcutil.Address `json:"address"`
	Amount  int64           `json:"amount"`
}

type BtcTransaction struct {
	txauthor.AuthoredTx
	chainParams *chaincfg.Params
}

func NewBtcTransaction(unspents []BtcUnspent, outputs []BtcOutput,
	changeAddress btcutil.Address, feeRate int, chainCfg *chaincfg.Params) (*BtcTransaction, error) {

	if len(unspents) == 0 || changeAddress == nil || feeRate <= 0 {
		return nil, errors.New("wrong params")
	}
	if !changeAddress.IsForNet(chainCfg) {
		return nil, errors.New("change address is not the corresponding network address")
	}
	changeBytes, err := txscript.PayToAddrScript(changeAddress)
	if err != nil {
		return nil, err
	}

	feeRatePerKb := btcutil.Amount(int64(feeRate) * 1000)

	txOuts, err := makeTxOutputs(outputs, feeRatePerKb, chainCfg)
	if err != nil {
		return nil, err
	}

	changeSource := txauthor.ChangeSource{
		NewScript: func() ([]byte, error) {
			return changeBytes, nil
		},
		ScriptSize: len(changeBytes),
	}

	unsignedTx, err := txauthor.NewUnsignedTransaction(txOuts, feeRatePerKb, makeInputSource(unspents), &changeSource)
	if err != nil {
		return nil, err
	}
	// Randomize change position, if change exists, before signing.  This
	// doesn't affect the serialize size, so the change amount will still
	// be valid.
	if unsignedTx.ChangeIndex >= 0 {
		unsignedTx.RandomizeChangePosition()
	}

	return &BtcTransaction{*unsignedTx, chainCfg}, nil
}

func (t *BtcTransaction) Sign(wallet *wallet.BtcWallet) error {
	return t.SignWithSecretsSource(wallet)
}

func (t *BtcTransaction) SignWithSecretsSource(secretsSource txauthor.SecretsSource) error {
	err := t.AddAllInputScripts(secretsSource)
	if err != nil {
		return err
	}
	err = t.validate()
	if err != nil {
		return err
	}

	return nil
}

func (t *BtcTransaction) GetFee() int64 {
	fee := t.TotalInput - txauthor.SumOutputValues(t.Tx.TxOut)
	return int64(fee)
}

func (t *BtcTransaction) Serialize() (string, error) {
	// Serialize the transaction and convert to hex string.
	buf := bytes.NewBuffer(make([]byte, 0, t.Tx.SerializeSize()))
	if err := t.Tx.Serialize(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (t *BtcTransaction) GetTxid() string {
	return t.Tx.TxHash().String()
}

func (t *BtcTransaction) Decode() *btcjson.TxRawDecodeResult {
	return DecodeMsgTx(t.Tx, t.chainParams)
}

func (t *BtcTransaction) Send(client *rpcclient.Client, allowHighFees bool) (*chainhash.Hash, error) {
	hash, err := client.SendRawTransaction(t.Tx, allowHighFees)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func (t *BtcTransaction) validate() error {
	hashCache := txscript.NewTxSigHashes(t.Tx)
	for i, prevScript := range t.PrevScripts {
		vm, err := txscript.NewEngine(prevScript, t.Tx, i,
			txscript.StandardVerifyFlags, nil, hashCache, int64(t.PrevInputValues[i]))
		if err != nil {
			return fmt.Errorf("cannot create script engine: %s", err)
		}
		err = vm.Execute()
		if err != nil {
			return fmt.Errorf("cannot validate transaction: %s", err)
		}
	}
	return nil
}

func makeTxOutputs(outputs []BtcOutput, relayFeePerKb btcutil.Amount, chainCfg *chaincfg.Params) ([]*wire.TxOut, error) {
	outLen := len(outputs)
	if outLen == 0 {
		return nil, errors.New("tx output is empty")
	}

	txOuts := make([]*wire.TxOut, 0, outLen)
	for i := 0; i < outLen; i++ {
		out := &outputs[i]

		if !out.Address.IsForNet(chainCfg) {
			return nil, errors.New("out address is not the corresponding network address")
		}

		// Create a new script which pays to the provided address.
		pkScript, err := txscript.PayToAddrScript(out.Address)
		if err != nil {
			return nil, err
		}
		txOut := &wire.TxOut{
			Value:    out.Amount,
			PkScript: pkScript,
		}
		if err = txrules.CheckOutput(txOut, relayFeePerKb); err != nil {
			return nil, err
		}

		txOuts = append(txOuts, txOut)
	}
	return txOuts, nil
}

func makeInputSource(unspents []BtcUnspent) txauthor.InputSource {
	sz := len(unspents)
	// Current inputs and their total value.  These are closed over by the
	// returned input source and reused across multiple calls.
	currentTotal := btcutil.Amount(0)
	currentInputs := make([]*wire.TxIn, 0, sz)
	currentInputValues := make([]btcutil.Amount, 0, sz)
	currentScripts := make([][]byte, 0, sz)

	return func(target btcutil.Amount) (btcutil.Amount, []*wire.TxIn, []btcutil.Amount, [][]byte, error) {
		for currentTotal < target && len(unspents) != 0 {
			u := unspents[0]
			unspents = unspents[1:]

			hash, _ := chainhash.NewHashFromStr(u.TxID)
			nextInput := wire.NewTxIn(&wire.OutPoint{
				Hash:  *hash,
				Index: u.Vout,
			}, nil, nil)

			amount, _ := btcutil.NewAmount(u.Amount)
			s, _ := hex.DecodeString(u.ScriptPubKey)

			currentTotal += amount
			currentInputs = append(currentInputs, nextInput)
			currentInputValues = append(currentInputValues, amount)
			currentScripts = append(currentScripts, s)
		}
		return currentTotal, currentInputs, currentInputValues, currentScripts, nil
	}
}
