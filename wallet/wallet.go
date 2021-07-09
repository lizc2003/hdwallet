package wallet

type Wallet interface {
	ChainId() int
	Symbol() string

	DeriveAddress() string
	DerivePublicKey() string
	DerivePrivateKey() string
}
