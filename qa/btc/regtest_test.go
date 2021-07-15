package btc

import (
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcutil"
	"github.com/lizc2003/hdwallet/btc"
	"github.com/lizc2003/hdwallet/wallet"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegtest(t *testing.T) {
	rq := require.New(t)

	cli, killBitcoind, err := RunBitcoind(&RunOptions{NewTmpDir: true})
	rq.Nil(err)
	defer killBitcoind()

	mnemonic, err := wallet.NewMnemonic(128)
	rq.Nil(err)

	btcChainId := wallet.BtcChainRegtest
	hdw, err := wallet.NewHDWallet(mnemonic, "", btcChainId, wallet.ChainMainNet)
	rq.Nil(err)

	w0, err := hdw.NewWallet(wallet.SymbolBtc, 0, 0, 0)
	//w0, err := hdw.NewSegWitWallet(0, 0, 0)
	//w0, err := hdw.NewNativeSegWitWallet(0, 0, 0)
	rq.Nil(err)

	w1, err := hdw.NewWallet(wallet.SymbolBtc, 0, 0, 1)
	rq.Nil(err)
	w2, err := hdw.NewSegWitWallet(0, 0, 1)
	rq.Nil(err)
	w3, err := hdw.NewNativeSegWitWallet(0, 0, 1)
	rq.Nil(err)

	chainCfg, _ := wallet.GetBtcChainConfig(btcChainId)
	a0 := w0.DeriveAddress()
	a1 := w1.DeriveAddress()
	a2 := w2.DeriveAddress()
	a3 := w3.DeriveAddress()
	fmt.Printf("a0: %s\na1: %s\na2: %s\na3: %s\n", a0, a1, a2, a3)
	addrA0, _ := btc.DecodeAddress(a0, chainCfg)
	addrA1, _ := btc.DecodeAddress(a1, chainCfg)
	addrA2, _ := btc.DecodeAddress(a2, chainCfg)
	addrA3, _ := btc.DecodeAddress(a3, chainCfg)

	{
		for _, ad := range []string{a0, a1, a2, a3} {
			err = cli.RpcClient.ImportAddress(ad)
			rq.Nil(err)
		}
	}

	var utxo btcjson.ListUnspentResult
	{
		// Generate 101 blocks for a0 (produce utxo)
		_, err = cli.RpcClient.GenerateToAddress(101, addrA0, nil)
		rq.Nil(err)

		cliUnspents, err := cli.RpcClient.ListUnspentMinMaxAddresses(1, 999, []btcutil.Address{addrA0})
		rq.Nil(err)
		rq.Equal(1, len(cliUnspents), "")

		utxo = cliUnspents[0]
		rq.Equal(50.0, utxo.Amount, "coinbase amount should be 50")
		fmt.Println(utxo)
	}

	transferAmount := 2.2
	var tx *btc.BtcTransaction

	{ // build tx
		unspent := btc.BtcUnspent{TxID: utxo.TxID, Vout: utxo.Vout,
			ScriptPubKey: utxo.ScriptPubKey, RedeemScript: utxo.RedeemScript,
			Amount: utxo.Amount}

		out1 := btc.BtcOutput{Address: addrA1, Amount: btc.BtcToSatoshi(transferAmount)}
		out2 := btc.BtcOutput{Address: addrA2, Amount: btc.BtcToSatoshi(transferAmount)}
		out3 := btc.BtcOutput{Address: addrA3, Amount: btc.BtcToSatoshi(transferAmount)}
		feeRate := 80

		tx, err = btc.NewBtcTransaction([]btc.BtcUnspent{unspent}, []btc.BtcOutput{out1, out2, out3},
			addrA0, feeRate, chainCfg)
		rq.Nil(err)
	}

	{ // fee
		fee := tx.GetFee()
		fmt.Println("fee:", fee)
	}

	{ // sign
		err = tx.Sign(w0.(*wallet.BtcWallet))
		rq.Nil(err)
	}

	{ // decode
		ret := tx.Decode()
		b, _ := json.MarshalIndent(ret, "", " ")
		fmt.Println("decoded tx:", string(b))
	}

	{ // send
		hash, err := tx.Send(cli.RpcClient, false)
		rq.Nil(err)
		txid := hash.String()
		fmt.Println("txid:", txid)
		rq.Equal(txid, tx.GetTxid(), "txid mismatch")

		rawtx, err := cli.RpcClient.GetRawTransactionVerbose(hash)
		rq.Nil(err)
		b, _ := json.MarshalIndent(rawtx, "", " ")
		fmt.Println("raw tx:", string(b))

		hex, _ := tx.Serialize()
		rq.Equal(hex, rawtx.Hex, "wrong hex")
	}

	{ //generate 1 block
		_, err = cli.RpcClient.GenerateToAddress(1, addrA0, nil)
		rq.Nil(err)
	}

	{ // query utxo of a1
		utxos, err := cli.RpcClient.ListUnspentMinMaxAddresses(0, 999, []btcutil.Address{addrA1})
		rq.Nil(err)
		rq.True(len(utxos) > 0, "No utxo of a1 fund!")
		b, _ := json.MarshalIndent(utxos, "", " ")
		fmt.Println("utxo for a1:", string(b))
		rq.Equal(transferAmount, utxos[0].Amount, "Wrong amount")
	}
}
