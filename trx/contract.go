package trx

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/contract"
	"github.com/lizc2003/hdwallet/wallet"
)

func DeployContract(w *wallet.TrxWallet, client *client.GrpcClient, contractName string,
	jsonAbi string, bytecode string,
	feeLimit, consumeUserResourcePercent, originEnergyLimit int64) (string, error) {

	abi, err := contract.JSONtoABI(jsonAbi)
	if err != nil {
		return "", err
	}
	txExt, err := client.DeployContract(w.DeriveAddress(), contractName, abi, bytecode,
		feeLimit, consumeUserResourcePercent, originEnergyLimit)
	if err != nil {
		return "", err
	}

	return SignAndSendTx(w, client, txExt)
}

func GetContractAddress(client *client.GrpcClient, txId string) (string, error) {
	info, err := client.GetTransactionInfoByID(txId)
	if err != nil {
		return "", err
	}

	return EncodeAddress(info.ContractAddress), nil
}
