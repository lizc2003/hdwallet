package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"
)

type TrxWallet struct {
	symbol     string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewTrxWallet(privateKey string) (*TrxWallet, error) {
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := derivePublicKey(privKey)
	if err != nil {
		return nil, err
	}

	return &TrxWallet{symbol: SymbolTrx,
		privateKey: privKey, publicKey: publicKey}, nil
}

func NewTrxWalletByPath(path string, seed []byte) (*TrxWallet, error) {
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	privKey, err := DerivePrivateKeyByPath(masterKey, path, false)
	if err != nil {
		return nil, err
	}
	privateKey := privKey.ToECDSA()

	publicKey, err := derivePublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	return &TrxWallet{symbol: SymbolTrx,
		privateKey: privateKey, publicKey: publicKey}, nil
}

func (w *TrxWallet) ChainId() int {
	return 0
}

func (w *TrxWallet) Symbol() string {
	return w.symbol
}

func (w *TrxWallet) DeriveAddress() string {
	const addressPrefix = 0x41
	return base58.CheckEncode(crypto.PubkeyToAddress(*w.publicKey).Bytes(), addressPrefix)
}

func (w *TrxWallet) DerivePublicKey() string {
	return hex.EncodeToString(crypto.FromECDSAPub(w.publicKey))
}

func (w *TrxWallet) DerivePrivateKey() string {
	return hex.EncodeToString(crypto.FromECDSA(w.privateKey))
}
