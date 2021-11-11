package eth

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type EthClient struct {
	RpcClient *ethclient.Client
	client    *rpc.Client
}

func NewEthClient(URL string) (*EthClient, error) {
	client, err := rpc.Dial(URL)
	if err != nil {
		return nil, err
	}
	rpcClient := ethclient.NewClient(client)
	return &EthClient{RpcClient: rpcClient, client: client}, nil
}

func (this *EthClient) SetHeader(key, value string) {
	if this.client != nil {
		this.client.SetHeader(key, value)
	}
}

func (this *EthClient) GetTransactionCountByNumber(ctx context.Context, blockNumber int64) (uint, error) {
	var num hexutil.Uint
	err := this.client.CallContext(ctx, &num, "eth_getBlockTransactionCountByNumber", hexutil.EncodeBig(big.NewInt(blockNumber)))
	return uint(num), err
}

func (this *EthClient) SendRawTransaction(ctx context.Context, signedHex string) (string, error) {
	var txid string
	err := this.client.CallContext(ctx, &txid, "eth_sendRawTransaction", signedHex)
	if err == nil && txid == "" {
		err = errors.New("SendRawTransaction: txid is empty")
	}
	return txid, err
}
