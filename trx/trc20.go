package trx

import (
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/lizc2003/hdwallet/wallet"
	"math/big"
)

type Trc20Contract struct {
	contractAddress string
	client          *client.GrpcClient
}

func NewErc20Contract(address string, client *client.GrpcClient) *Trc20Contract {
	return &Trc20Contract{contractAddress: address, client: client}
}

func (this *Trc20Contract) Name() (string, error) {
	return this.client.TRC20GetName(this.contractAddress)
}

func (this *Trc20Contract) Symbol() (string, error) {
	return this.client.TRC20GetSymbol(this.contractAddress)
}

func (this *Trc20Contract) Decimals() (int, error) {
	n, err := this.client.TRC20GetDecimals(this.contractAddress)
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}

func (this *Trc20Contract) BalanceOf(tokenOwner string) (*big.Int, error) {
	return this.client.TRC20ContractBalance(tokenOwner, this.contractAddress)
}

func (this *Trc20Contract) Transfer(w *wallet.TrxWallet, toAddr string, amount *big.Int, feeLimit int64) (string, error) {
	txExt, err := this.client.TRC20Send(w.DeriveAddress(), toAddr, this.contractAddress, amount, feeLimit)
	if err != nil {
		return "", err
	}

	return SignAndSendTx(w, this.client, txExt)
}
