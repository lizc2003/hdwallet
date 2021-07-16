package btc

import (
	"errors"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/lizc2003/hdwallet/wallet"
)

type BtcClient struct {
	RpcClient *rpcclient.Client
}

func NewBtcClient(host string, user string, pass string, chainId int) (*BtcClient, error) {
	chainCfg, err := wallet.GetBtcChainConfig(chainId)
	if err != nil {
		return nil, err
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         pass,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
		Params:       chainCfg.Name,
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		return nil, err
	}

	return &BtcClient{RpcClient: client}, nil
}

func (this *BtcClient) EstimateFeePerKb() (int64, error) {
	feeResult, err := this.RpcClient.EstimateSmartFee(6, nil)
	if err != nil {
		return 0, err
	}

	if feeResult.FeeRate != nil && *feeResult.FeeRate > 0 {
		return BtcToSatoshi(*feeResult.FeeRate), nil
	} else {
		return 0, errors.New("Fee not available")
	}
}
