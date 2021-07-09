package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

type EthClient struct {
	RpcClient *ethclient.Client
	client    *rpc.Client
}

func NewEthClient(rawurl string) (*EthClient, error) {
	client, err := rpc.Dial(rawurl)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(client)
	return &EthClient{RpcClient: ethClient, client: client}, nil
}

func (this *EthClient) GetTransactionCountByNumber(ctx context.Context, blockNumber int64) (uint, error) {
	var num hexutil.Uint
	err := this.client.CallContext(ctx, &num, "eth_getBlockTransactionCountByNumber", hexutil.EncodeBig(big.NewInt(blockNumber)))
	return uint(num), err
}
