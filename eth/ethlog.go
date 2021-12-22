package eth

import (
	"context"
	"encoding/hex"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"strings"
)

const (
	TokenTypeErc20  = 1
	TokenTypeErc721 = 2
)

type TransferEvent struct {
	Txid            string
	BlockNumber     int64
	TokenType       int
	ContractAddress string
	From            string
	To              string
	Value           string
}

func FilterTransferLog(cli *ethclient.Client, fromBlock int64, toBlock int64, contractAddresses []common.Address) ([]TransferEvent, error) {
	//hex.EncodeToString(crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).Bytes()))
	transferTopic := "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	logs, err := cli.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: contractAddresses,
		Topics:    [][]common.Hash{[]common.Hash{HexToHash(transferTopic)}},
	})
	if err != nil {
		return nil, err
	}

	var evts []TransferEvent
	sz := len(logs)
	if sz > 0 {
		evts = make([]TransferEvent, 0, sz)
		for i := 0; i < sz; i++ {
			log := &logs[i]
			topicLen := len(log.Topics)
			if topicLen == 3 || topicLen == 4 {
				if hex.EncodeToString(log.Topics[0].Bytes()) != transferTopic { // Transfer event
					continue
				}

				from := parseEthLogAddress(hex.EncodeToString(log.Topics[1].Bytes()))
				to := parseEthLogAddress(hex.EncodeToString(log.Topics[2].Bytes()))

				var val string
				var tokenType int
				if topicLen == 3 { // ERC20 Transfer event
					tokenType = TokenTypeErc20
					val = hex.EncodeToString(log.Data)
					if len(val) == 64 {
						val = parseEthLogValue(val)
					}
				} else if topicLen == 4 { // ERC721 Transfer event
					tokenType = TokenTypeErc721
					val = parseEthLogValue(hex.EncodeToString(log.Topics[3].Bytes()))
				}
				if val != "" {
					evt := TransferEvent{
						Txid:            log.TxHash.String(),
						BlockNumber:     int64(log.BlockNumber),
						TokenType:       tokenType,
						ContractAddress: strings.ToLower(log.Address.String()),
						From:            from, To: to, Value: val,
					}
					evts = append(evts, evt)
				}
			}
		}
	}
	return evts, nil
}

func parseEthLogAddress(addr string) string {
	sz := len(addr)
	if sz >= 40 {
		return "0x" + strings.ToLower(addr[sz-40:])
	} else {
		return ""
	}
}

func parseEthLogValue(val string) string {
	value := strings.TrimLeft(val, "0")
	sz := len(value)
	if sz == 0 {
		value = "00"
	} else if sz%2 == 1 {
		value = "0" + value
	}
	return value
}
