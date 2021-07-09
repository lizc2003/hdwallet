package eth

type Erc20Meta struct {
	Symbol   string `toml:"symbol" json:"symbol"`
	Address  string `toml:"address" json:"address"`
	Decimals int    `toml:"decimals" json:"decimals"`
}
