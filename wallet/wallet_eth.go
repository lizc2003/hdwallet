package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

type EthWallet struct {
	symbol     string
	chainId    int
	chainCfg   *params.ChainConfig
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func NewEthWallet(privateKey string, chainId int) (*EthWallet, error) {
	chainCfg, err := GetEthChainConfig(chainId)
	if err != nil {
		return nil, err
	}

	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	publicKey, err := derivePublicKey(privKey)
	if err != nil {
		return nil, err
	}

	return &EthWallet{symbol: SymbolEth,
		chainId: chainId, chainCfg: chainCfg,
		privateKey: privKey, publicKey: publicKey}, nil
}

func NewEthWalletByPath(path string, seed []byte, chainId int) (*EthWallet, error) {
	chainCfg, err := GetEthChainConfig(chainId)
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

	publicKey, err := derivePublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	return &EthWallet{symbol: SymbolEth,
		chainId: chainId, chainCfg: chainCfg,
		privateKey: privateKey, publicKey: publicKey}, nil
}

func (w *EthWallet) ChainId() int {
	return w.chainId
}

func (w *EthWallet) ChainParams() *params.ChainConfig {
	return w.chainCfg
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

func derivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to get public key")
	}

	return publicKeyECDSA, nil
}

func (w *EthWallet) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	signer := types.LatestSigner(w.chainCfg)
	signedTx, err := types.SignTx(tx, signer, w.privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}
