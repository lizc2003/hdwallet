package wallet

import (
	"fmt"
	"strings"
)

type HDWallet struct {
	seed       []byte
	btcChainId int
	ethChainId int
}

func NewHDWallet(mnemonic, password string, btcChainId int, ethChainId int) (*HDWallet, error) {
	mnemonic = strings.ReplaceAll(mnemonic, "\n", "")
	mnemonic = strings.ReplaceAll(mnemonic, "\r", "")

	seed, err := NewSeedFromMnemonic(mnemonic, password)
	if err != nil {
		return nil, err
	}
	return &HDWallet{seed: seed, btcChainId: btcChainId, ethChainId: ethChainId}, nil
}

func (this *HDWallet) NewWallet(symbol string, accountIndex, changeType, index int) (Wallet, error) {
	path, err := MakeBip44Path(symbol, this.btcChainId, accountIndex, changeType, index)
	if err != nil {
		return nil, err
	}

	return this.NewWalletByPath(symbol, path, SegWitNone)
}

func (this *HDWallet) NewSegWitWallet(accountIndex, changeType, index int) (Wallet, error) {
	path, err := MakeBip49Path(SymbolBtc, this.btcChainId, accountIndex, changeType, index)
	if err != nil {
		return nil, err
	}
	return this.NewWalletByPath(SymbolBtc, path, SegWitScript)
}

func (this *HDWallet) NewNativeSegWitWallet(accountIndex, changeType, index int) (Wallet, error) {
	path, err := MakeBip84Path(SymbolBtc, this.btcChainId, accountIndex, changeType, index)
	if err != nil {
		return nil, err
	}
	return this.NewWalletByPath(SymbolBtc, path, SegWitNative)
}

func (this *HDWallet) NewWalletByPath(symbol string, path string, segWitType SegWitType) (Wallet, error) {
	var w Wallet
	var err error

	switch symbol {
	case SymbolBtc:
		w, err = NewBtcWalletByPath(path, this.seed, this.btcChainId, segWitType)
	case SymbolEth:
		w, err = NewEthWalletByPath(path, this.seed, this.ethChainId)
	default:
		err = fmt.Errorf("invalid symbol: %s", symbol)
	}

	if err != nil {
		return nil, err
	}
	return w, nil
}
