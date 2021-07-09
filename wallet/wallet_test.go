package wallet

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCoin_MakeAddress(t *testing.T) {
	mnemonic, err := NewMnemonic(256)
	require.NoError(t, err)
	t.Log(mnemonic)

	btcChainId := BtcChainMainNet
	ethChainId := ChainMainNet

	hdw, err := NewHDWallet(mnemonic, "", btcChainId, ethChainId)
	require.NoError(t, err)

	w, err := hdw.NewWallet(SymbolBtc, 0, 0, 0)
	require.NoError(t, err)

	privateKey := w.DerivePrivateKey()
	publicKey := w.DerivePublicKey()
	address := w.DeriveAddress()
	t.Log(privateKey)
	t.Log(publicKey)
	t.Log(address)

	w2, err := NewBtcWallet(privateKey, btcChainId, SegWitNone)
	require.NoError(t, err)
	privateKey2 := w2.DerivePrivateKey()
	publicKey2 := w2.DerivePublicKey()
	address2 := w2.DeriveAddress()
	require.Equal(t, privateKey, privateKey2)
	require.Equal(t, publicKey, publicKey2)
	require.Equal(t, address, address2)

	w, err = hdw.NewWallet(SymbolEth, 0, 0, 0)
	require.NoError(t, err)

	privateKey = w.DerivePrivateKey()
	publicKey = w.DerivePublicKey()
	address = w.DeriveAddress()
	t.Log(privateKey)
	t.Log(publicKey)
	t.Log(address)

	w3, err := NewEthWallet(privateKey, ethChainId)
	require.NoError(t, err)
	privateKey3 := w3.DerivePrivateKey()
	publicKey3 := w3.DerivePublicKey()
	address3 := w3.DeriveAddress()
	require.Equal(t, privateKey, privateKey3)
	require.Equal(t, publicKey, publicKey3)
	require.Equal(t, address, address3)
}

func TestCoin_MakeAddressTestnet(t *testing.T) {
	mnemonic := "purse cheese cage reason cost flat jump usage hospital grit delay loan"
	hdw, err := NewHDWallet(mnemonic, "", BtcChainTestNet3, ChainGoerli)
	require.NoError(t, err)

	w, err := hdw.NewWallet(SymbolBtc, 0, 0, 0)
	require.NoError(t, err)

	privateKey := w.DerivePrivateKey()
	publicKey := w.DerivePublicKey()
	address := w.DeriveAddress()
	require.Equal(t, privateKey, "cUpt1n6hoxtYWGNX4jiTGA31nroMgFXFsYzmHs8cvDfCnGLjH5MV")
	require.Equal(t, publicKey, "03e8a1810e3063ccd7a3e1f40195e13dcb0678e0097433715e9cf8e31fce753b10")
	require.Equal(t, address, "mt1HJYiAki1FDfGh8M7MEeoKc8BKJ4hhc6")

	w, err = hdw.NewWallet(SymbolBtc, 0, 0, 1)
	require.NoError(t, err)

	privateKey = w.DerivePrivateKey()
	publicKey = w.DerivePublicKey()
	address = w.DeriveAddress()
	require.Equal(t, privateKey, "cS2B6pjUu6gdwWCHuPdFuCz1xWh47TGqnqXnvCAu8vfbf9a6fyGv")
	require.Equal(t, publicKey, "02d5a38f5f39dff2c8b57b0950acc285822cc1d9efca6c80abe1f0a92f90bfcc1b")
	require.Equal(t, address, "n2BVx6LNcFkSoqxc7EatFAByoZXwmvztxJ")

	w, err = hdw.NewWallet(SymbolEth, 0, 0, 0)
	require.NoError(t, err)

	privateKey = w.DerivePrivateKey()
	publicKey = w.DerivePublicKey()
	address = w.DeriveAddress()
	require.Equal(t, privateKey, "e16ac20fafb7de15445488f1fc6a0e5a05e9efca52acb15de559e4914c8f351d")
	require.Equal(t, publicKey, "040861286e63b01682f3769628934292819d37dbdacb48994ebde5d54ddc4146dccf61b578513d0a27c1970984e6def2d31832ff5f5df7b16fe16bf4b2b87d003c")
	require.Equal(t, address, "0x942429aA212ef7cb14DfE06ed0EE78EB82BD298f")

	w, err = hdw.NewWallet(SymbolEth, 0, 0, 1)
	require.NoError(t, err)

	privateKey = w.DerivePrivateKey()
	publicKey = w.DerivePublicKey()
	address = w.DeriveAddress()
	require.Equal(t, privateKey, "63191752898f7a7f20caa49da570d67e4fc0c2623a92c2fe3401d427a2de77aa")
	require.Equal(t, publicKey, "042740014857de4394a29515a7a80f192e5a335440d20ad73d662901e7f82720f6060b63b470a5a31ec45deba053b259966e3fed3c62b0fd6371abb4120d53555e")
	require.Equal(t, address, "0x3295Db1E775723c752511c4E4caA400dbaf2240F")
}

func TestCoin_SegWitAddress(t *testing.T) {
	mnemonic := "example escape erode educate help cigar super chalk best inner fossil soft"
	hdw, err := NewHDWallet(mnemonic, "", BtcChainMainNet, ChainMainNet)
	require.NoError(t, err)

	w, err := hdw.NewSegWitWallet(0, 0, 0)
	require.NoError(t, err)
	privateKey := w.DerivePrivateKey()
	publicKey := w.DerivePublicKey()
	witAddress := w.DeriveAddress()
	t.Log(privateKey)
	t.Log(publicKey)
	t.Log(witAddress)
	require.Equal(t, witAddress, "32dBJzEpeCnuk5ti8Pa2o232KYHV2WGmkc")

	w, err = hdw.NewNativeSegWitWallet(0, 0, 0)
	require.NoError(t, err)
	privateKey = w.DerivePrivateKey()
	publicKey = w.DerivePublicKey()
	witAddress = w.DeriveAddress()
	t.Log(" ")
	t.Log(privateKey)
	t.Log(publicKey)
	t.Log(witAddress)
	require.Equal(t, witAddress, "bc1qtmygf2u2n3wdgepgtfxvzy4xr83gzc9cyqf594")
}

func TestCoin_DeriveAddress(t *testing.T) {
	const mnemonic = "eternal list thank chaos trick paper sniff ridge make govern invest abandon"

	hdw, err := NewHDWallet(mnemonic, "lixu1234qwer", BtcChainTestNet3, ChainGoerli)
	require.NoError(t, err)

	for i, tt := range []struct {
		privateKey string
		publicKey  string
		address    string
	}{
		{
			privateKey: "cVfZerHznFjwSC85exvS1B9xpqU2yb6piLHQmUN5DnsHek3JV7xn",
			publicKey:  "031cf3493c5fcb4eabdfaa4191a02cc30429539ea6b80f5590bc4a8b6222f0d3ba",
			address:    "mm16s7xsf8Wjwxhprc6YzLW9gVncqZNGBR",
		},
		{
			privateKey: "cMf15hojxybmvPfkBD3pGVKKHj4anSqD4SxuH9egUxKmhGDPBcLh",
			publicKey:  "02d3c0c2de32c8923be88c45aa56bf30aebd428feb13e578e946c64668425b82ee",
			address:    "n23dZRLMrWF419JNe28WWCTatcNBe3PCA9",
		},
		{
			privateKey: "cPpmQ8iDoj142QwQb2xBuK2ScCR6N2LKQBEATbSn5pTrpaKnZ6TB",
			publicKey:  "03bf34a6d18b36459432644fd3a41a9b5ae985c95ee8627613366ec7bd13ea0af4",
			address:    "mzxzKR9Qz3zmQmRpKChRuDAS6WSiiuscGM",
		},
	} {
		w, err := hdw.NewWallet(SymbolBtc, 0, 0, i)
		require.NoError(t, err)

		privateKey := w.DerivePrivateKey()
		publicKey := w.DerivePublicKey()
		address := w.DeriveAddress()
		require.Equal(t, privateKey, tt.privateKey)
		require.Equal(t, publicKey, tt.publicKey)
		require.Equal(t, address, tt.address)
	}

	hdw, err = NewHDWallet(mnemonic, "lixu1234qwer", BtcChainMainNet, ChainMainNet)
	require.NoError(t, err)

	for i, tt := range []struct {
		privateKey string
		publicKey  string
		address    string
	}{
		{
			privateKey: "L1evYPBu2pWjWCKkKmy5aYcvNSx1FYqhZjJDepQTVmE38RAFAboN",
			publicKey:  "0209f66c7e3f0c42c82dcc17fd6ccd51b3a01ec55d2c3dadd0b7c5456dcaee97ff",
			address:    "1Gh3TUC3WLhB4t6Km9fSPxJLVX4dUBFiVD",
		},
		{
			privateKey: "L5EhohWQ2rm5a1aF27akRtfZTXfmHVt2Vhxrs1C6s8RS6JonVr1X",
			publicKey:  "0263b1066c8626630e19c9d3ccd115fc808a791567b2b6abcdf9864fb62aaca9ae",
			address:    "1H7JDuCF8QBcy28GUgNvB9wscdjw1ZrrPc",
		},
		{
			privateKey: "KwX6JPubzTxQRLLS6tVDjX2aAwE4nDAqSGJoaPxqHf3EowybHRbd",
			publicKey:  "034968dc4f6e35ff7d478edcd17af1eaedfda97bcc39545129e725e5c6f164c343",
			address:    "12AxjDAM6fhPsf66gEiFMYcY3uSdSUi9Zj",
		},
	} {
		w, err := hdw.NewWallet(SymbolBtc, 0, 0, i)
		require.NoError(t, err)

		privateKey := w.DerivePrivateKey()
		publicKey := w.DerivePublicKey()
		address := w.DeriveAddress()
		require.Equal(t, privateKey, tt.privateKey)
		require.Equal(t, publicKey, tt.publicKey)
		require.Equal(t, address, tt.address)
	}
}

func TestNewMnemonic(t *testing.T) {
	mn, err := NewMnemonic(128)
	assert.NoError(t, err)
	en, err := EntropyFromMnemonic(mn)
	assert.NoError(t, err)
	mn1, err := NewMnemonicByEntropy(en)
	assert.NoError(t, err)
	assert.EqualValues(t, mn, mn1)
}
