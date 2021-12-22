package eth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"strconv"
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

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

func (this *EthClient) TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, from *common.Address, blockNumber *big.Int, err error) {
	var json *rpcTransaction
	err = this.client.CallContext(ctx, &json, "eth_getTransactionByHash", hash)
	if err != nil {
		return nil, nil, nil, err
	} else if json == nil {
		return nil, nil, nil, ethereum.NotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, nil, nil, fmt.Errorf("server returned transaction without signature")
	}
	if json.BlockNumber != nil {
		if tmp, err := strconv.ParseInt(*json.BlockNumber, 0, 64); err == nil {
			blockNumber = big.NewInt(tmp)
		}
	}
	return json.tx, json.From, blockNumber, nil
}
