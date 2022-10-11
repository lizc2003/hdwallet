package btc

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
)

func DecodeAddress(addr string, chainParams *chaincfg.Params) (btcutil.Address, error) {
	return btcutil.DecodeAddress(addr, chainParams)
}

func HexToHash(s string) (*chainhash.Hash, error) {
	return chainhash.NewHashFromStr(s)
}

func BtcToSatoshi(v float64) int64 {
	amt, _ := btcutil.NewAmount(v)
	return int64(amt)
}

func SatoshiToBtc(v int64) float64 {
	a := btcutil.Amount(v)
	return a.ToBTC()
}

func EstimateFee(numP2PKHIns, numP2WPKHIns, numNestedP2WPKHIns int,
	outputs []BtcOutput, feePerKb int64, changeScriptSize int, chainCfg *chaincfg.Params) (int64, int64, error) {

	feeRatePerKb := btcutil.Amount(feePerKb)
	if changeScriptSize < 0 {
		// using P2WPKH as change output.
		changeScriptSize = txsizes.P2WPKHPkScriptSize
	}

	txOuts, err := makeTxOutputs(outputs, feeRatePerKb, chainCfg)
	if err != nil {
		return 0, 0, err
	}

	maxSignedSize := txsizes.EstimateVirtualSize(numP2PKHIns, 0, numP2WPKHIns,
		numNestedP2WPKHIns, txOuts, changeScriptSize)

	targetFee := txrules.FeeForSerializeSize(feeRatePerKb, maxSignedSize)
	targetAmount := txauthor.SumOutputValues(txOuts)

	return int64(targetFee), int64(targetAmount), nil
}
