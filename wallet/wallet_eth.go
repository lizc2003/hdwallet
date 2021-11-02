package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type EthWallet struct {
	symbol      string
	chainId     int
	chainParams *params.ChainConfig
	privateKey  *ecdsa.PrivateKey
	publicKey   *ecdsa.PublicKey
}

func NewEthWallet(privateKey string, chainId int) (*EthWallet, error) {
	chainParams, err := GetEthChainParams(chainId)
	if err != nil {
		return nil, err
	}

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := DerivePublicKey(privKey)
	if err != nil {
		return nil, err
	}

	return &EthWallet{symbol: SymbolEth,
		chainId: chainId, chainParams: chainParams,
		privateKey: privKey, publicKey: publicKey}, nil
}

func NewEthWalletByPath(path string, seed []byte, chainId int) (*EthWallet, error) {
	chainParams, err := GetEthChainParams(chainId)
	if err != nil {
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	privKey, err := DerivePrivateKeyByPath(masterKey, path, IsFixIssue172)
	if err != nil {
		return nil, err
	}
	privateKey := privKey.ToECDSA()

	publicKey, err := DerivePublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	return &EthWallet{symbol: SymbolEth,
		chainId: chainId, chainParams: chainParams,
		privateKey: privateKey, publicKey: publicKey}, nil
}

func (w *EthWallet) ChainId() int {
	return w.chainId
}

func (w *EthWallet) ChainParams() *params.ChainConfig {
	return w.chainParams
}

func (w *EthWallet) Symbol() string {
	return w.symbol
}

func (w *EthWallet) DeriveAddress() string {
	return crypto.PubkeyToAddress(*w.publicKey).Hex()
}

func (w *EthWallet) DerivePublicKey() string {
	return hex.EncodeToString(crypto.FromECDSAPub(w.publicKey))
}

func (w *EthWallet) DerivePrivateKey() string {
	return hex.EncodeToString(crypto.FromECDSA(w.privateKey))
}

func (w *EthWallet) DeriveNativeAddress() common.Address {
	return crypto.PubkeyToAddress(*w.publicKey)
}

func (w *EthWallet) DeriveNativePrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}
