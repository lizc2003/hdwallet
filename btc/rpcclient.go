package btc

import (
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/lizc2003/hdwallet/wallet"
	"net/url"
)

type BtcClient struct {
	RpcClient *rpcclient.Client
}

func NewBtcClient(URL string, user string, pass string, chainId int) (*BtcClient, error) {
	chainCfg, err := wallet.GetBtcChainParams(chainId)
	if err != nil {
		return nil, err
	}

	connCfg := &rpcclient.ConnConfig{
		User:   user,
		Pass:   pass,
		Params: chainCfg.Name,
	}

	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		connCfg.HTTPPostMode = true
		connCfg.Host = u.Host + u.Path
		if u.Scheme == "http" {
			connCfg.DisableTLS = true
		}
	} else if u.Scheme == "ws" || u.Scheme == "wss" {
		connCfg.Host = u.Host
		if u.Path != "" {
			connCfg.Endpoint = u.Path[1:]
		}
		if u.Scheme == "ws" {
			connCfg.DisableTLS = true
		}
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

// https://bitcoincore.org/en/doc/0.21.0/rpc/rawtransactions/sendrawtransaction/
func (this *BtcClient) SendRawTransaction(signedHex string, allowHighFees bool) (string, error) {
	hex, _ := json.Marshal(signedHex)
	params := []json.RawMessage{hex}
	if allowHighFees {
		maxFeeRate, _ := json.Marshal(0)
		params = append(params, maxFeeRate)
	}

	resp, err := this.RpcClient.RawRequest("sendrawtransaction", params)
	if err != nil {
		return "", err
	}

	var txid string
	err = json.Unmarshal(resp, &txid)
	if err == nil && txid == "" {
		err = errors.New("unknown response")
	}
	if err != nil {
		return "", err
	}
	return txid, nil
}
