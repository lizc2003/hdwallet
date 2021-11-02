package wallet

import (
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/ethereum/go-ethereum/params"
	"github.com/tyler-smith/go-bip39"
	"math"
	"strconv"
)

type SegWitType int

const (
	SymbolEth = "ETH"
	SymbolBtc = "BTC"
	SymbolTrx = "TRX"

	BtcChainMainNet  = int(wire.MainNet)
	BtcChainTestNet3 = int(wire.TestNet3)
	BtcChainRegtest  = int(wire.TestNet)
	BtcChainSimNet   = int(wire.SimNet)

	ChainMainNet = 1    // for ETH
	ChainRopsten = 3    // for ETH
	ChainRinkeby = 4    // for ETH
	ChainGoerli  = 5    // for ETH
	ChainPrivate = 1337 // for ETH

	SegWitNone   SegWitType = 0
	SegWitScript SegWitType = 1
	SegWitNative SegWitType = 2

	ChangeTypeExternal = 0
	ChangeTypeInternal = 1 // Usually used for change, not visible to the outside world

	SatoshiPerBitcoin = 1e8
	GweiPerEther      = 1e9
	WeiPerGwei        = 1e9
	WeiPerEther       = 1e18

	EtherTransferGas = 21000

	TokenShowDecimals = 9
)

var IsFixIssue172 = false

func GetBtcChainParams(chainId int) (*chaincfg.Params, error) {
	switch chainId {
	case BtcChainMainNet:
		return &chaincfg.MainNetParams, nil
	case BtcChainTestNet3:
		return &chaincfg.TestNet3Params, nil
	case BtcChainRegtest:
		return &chaincfg.RegressionNetParams, nil
	case BtcChainSimNet:
		return &chaincfg.SimNetParams, nil
	default:
		return nil, fmt.Errorf("unknown btc chainId: %d", chainId)
	}
}

func GetEthChainParams(chainId int) (*params.ChainConfig, error) {
	switch chainId {
	case ChainMainNet:
		return params.MainnetChainConfig, nil
	case ChainRopsten:
		return params.RopstenChainConfig, nil
	case ChainRinkeby:
		return params.RinkebyChainConfig, nil
	case ChainGoerli:
		return params.GoerliChainConfig, nil
	case ChainPrivate:
		return params.TestChainConfig, nil
	default:
		return nil, fmt.Errorf("unknown eth chainId: %d", chainId)
	}
}

func NewEntropy(bits int) (entropy []byte, err error) {
	return bip39.NewEntropy(bits)
}

func NewMnemonic(bits int) (mnemonic string, err error) {
	entropy, err := NewEntropy(bits)
	if err != nil {
		return "", err
	}
	return NewMnemonicByEntropy(entropy)
}

func NewMnemonicByEntropy(entropy []byte) (mnemonic string, err error) {
	return bip39.NewMnemonic(entropy)
}

func EntropyFromMnemonic(mnemonic string) (entropy []byte, err error) {
	return bip39.EntropyFromMnemonic(mnemonic)
}

func NewSeedFromMnemonic(mnemonic, password string) ([]byte, error) {
	if mnemonic == "" {
		return nil, errors.New("mnemonic is required")
	}
	return bip39.NewSeedWithErrorChecking(mnemonic, password)
}

func MakeBip44Path(symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	return MakeBipXPath(44, symbol, chainId, accountIndex, changeType, index)
}

func MakeBip49Path(symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	return MakeBipXPath(49, symbol, chainId, accountIndex, changeType, index)
}

func MakeBip84Path(symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	return MakeBipXPath(84, symbol, chainId, accountIndex, changeType, index)
}

func MakeBipXPath(bipType int, symbol string, chainId int, accountIndex, changeType, index int) (string, error) {
	var coinType int
	switch symbol {
	case SymbolEth:
		coinType = 60
	case SymbolBtc:
		chainParams, err := GetBtcChainParams(chainId)
		if err != nil {
			return "", err
		}
		coinType = int(chainParams.HDCoinType)
	case SymbolTrx:
		coinType = 195
	default:
		return "", fmt.Errorf("invalid symbol: %s", symbol)
	}

	if accountIndex < 0 || index < 0 {
		return "", errors.New("invalid account index or index")
	}
	if changeType != ChangeTypeExternal && changeType != ChangeTypeInternal {
		return "", errors.New("invalid change type")
	}
	return fmt.Sprintf("m/%d'/%d'/%d'/%d/%d", bipType, coinType, accountIndex, changeType, index), nil
}

func FormatBtc(amount int64) string {
	return FormatFloat(float64(amount)/SatoshiPerBitcoin, 8)
}

func FormatEth(amount int64) string {
	return FormatFloat(float64(amount)/GweiPerEther, 9)
}

func FormatFloat(f float64, precision int) string {
	d := float64(1)
	if precision > 0 {
		d = math.Pow10(precision)
	}
	return strconv.FormatFloat(math.Trunc(f*d)/d, 'f', -1, 64)
}
