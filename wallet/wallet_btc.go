package wallet

import (
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"
	"log"
)

var ErrAddressNotMatch = errors.New("address not match")

type BtcWallet struct {
	symbol     string
	segWitType SegWitType
	chainCfg   *chaincfg.Params
	privateKey *btcec.PrivateKey
	publicKey  *btcec.PublicKey
}

func NewBtcWallet(privateKey string, chainId int, segWitType SegWitType) (*BtcWallet, error) {
	chainCfg, err := GetBtcChainConfig(chainId)
	if err != nil {
		return nil, err
	}

	wif, err := btcutil.DecodeWIF(privateKey)
	if err != nil {
		return nil, err
	}
	if !wif.IsForNet(chainCfg) {
		return nil, errors.New("key network doesn't match")
	}

	return &BtcWallet{symbol: SymbolBtc,
		chainCfg: chainCfg, segWitType: segWitType,
		privateKey: wif.PrivKey,
		publicKey:  wif.PrivKey.PubKey()}, nil
}

func NewBtcWalletByPath(path string, seed []byte, chainId int, segWitType SegWitType) (*BtcWallet, error) {
	chainCfg, err := GetBtcChainConfig(chainId)
	if err != nil {
		return nil, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, chainCfg)
	if err != nil {
		return nil, err
	}

	privateKey, err := DerivePrivateKeyByPath(masterKey, path, IsFixIssue172)
	if err != nil {
		return nil, err
	}

	return &BtcWallet{symbol: SymbolBtc,
		chainCfg: chainCfg, segWitType: segWitType,
		privateKey: privateKey,
		publicKey:  privateKey.PubKey()}, nil
}

func (w *BtcWallet) ChainId() int {
	return int(w.chainCfg.Net)
}

func (w *BtcWallet) ChainParams() *chaincfg.Params {
	return w.chainCfg
}

func (w *BtcWallet) Symbol() string {
	return w.symbol
}

func (w *BtcWallet) DeriveAddress() string {
	switch w.segWitType {
	case SegWitNone:
		pk := w.publicKey.SerializeCompressed()
		keyHash := btcutil.Hash160(pk)
		p2pkhAddr, err := btcutil.NewAddressPubKeyHash(keyHash, w.chainCfg)
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return ""
		}
		return p2pkhAddr.String()
	case SegWitScript:
		pk := w.publicKey.SerializeCompressed()
		keyHash := btcutil.Hash160(pk)
		scriptSig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_0).AddData(keyHash).Script()
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return ""
		}
		addr, err := btcutil.NewAddressScriptHash(scriptSig, w.chainCfg)
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return ""
		}
		return addr.String()
	case SegWitNative:
		pk := w.publicKey.SerializeCompressed()
		keyHash := btcutil.Hash160(pk)
		p2wpkh, err := btcutil.NewAddressWitnessPubKeyHash(keyHash, w.chainCfg)
		if err != nil {
			log.Println("DeriveAddress error:", err)
			return ""
		}
		return p2wpkh.String()
	}
	return ""
}

func (w *BtcWallet) DerivePublicKey() string {
	return hex.EncodeToString(w.publicKey.SerializeCompressed())
}

func (w *BtcWallet) DerivePrivateKey() string {
	wif, err := btcutil.NewWIF(w.privateKey, w.chainCfg, true)
	if err != nil {
		log.Println("DerivePrivateKey error:", err)
		return ""
	}
	return wif.String()
}

func DerivePrivateKeyByPath(masterKey *hdkeychain.ExtendedKey, path string, fixIssue172 bool) (*btcec.PrivateKey, error) {
	dpath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}

	key := masterKey
	for _, n := range dpath {
		if fixIssue172 && key.IsAffectedByIssue172() {
			key, err = key.Derive(n)
		} else {
			key, err = key.DeriveNonStandard(n)
		}
		if err != nil {
			return nil, err
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// txauthor.SecretsSource
func (w *BtcWallet) GetKey(addr btcutil.Address) (*btcec.PrivateKey, bool, error) {
	if w.DeriveAddress() == addr.EncodeAddress() {
		return w.privateKey, true, nil
	}
	return nil, false, ErrAddressNotMatch
}

func (w *BtcWallet) GetScript(addr btcutil.Address) ([]byte, error) {
	return nil, errors.New("GetScript not supported")
}
