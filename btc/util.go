package btc

import (
	"github.com/btcsuite/btcutil"
)

func BtcToSatoshi(v float64) int64 {
	amt, _ := btcutil.NewAmount(v)
	return int64(amt)
}

func SatoshiToBtc(v int64) float64 {
	a := btcutil.Amount(v)
	return a.ToBTC()
}
