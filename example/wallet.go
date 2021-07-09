package main

import (
	"fmt"
	"log"
	"github.com/lizc2003/hdwallet/wallet"
)

func main() {
	mnemonic, err := wallet.NewMnemonic(128)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mnemonic:", mnemonic)

	hdw, err := wallet.NewHDWallet(mnemonic, "", wallet.BtcChainMainNet, wallet.ChainMainNet)
	if err != nil {
		log.Fatal(err)
	}

	w, err := hdw.NewWallet(wallet.SymbolBtc, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wallet: %s\n\tprivatekey: %s\n\tpublickey: %s\n\taddress: %s\n",
		w.Symbol(),
		w.DerivePrivateKey(),
		w.DerivePublicKey(),
		w.DeriveAddress())

	w, err = hdw.NewSegWitWallet(0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("segwit wallet: %s\n\tprivatekey: %s\n\tpublickey: %s\n\taddress: %s\n",
		w.Symbol(),
		w.DerivePrivateKey(),
		w.DerivePublicKey(),
		w.DeriveAddress())

	w, err = hdw.NewNativeSegWitWallet(0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("native segwit wallet: %s\n\tprivatekey: %s\n\tpublickey: %s\n\taddress: %s\n",
		w.Symbol(),
		w.DerivePrivateKey(),
		w.DerivePublicKey(),
		w.DeriveAddress())

	w, err = hdw.NewWallet(wallet.SymbolEth, 0, 0, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wallet: %s\n\tprivatekey: %s\n\tpublickey: %s\n\taddress: %s\n",
		w.Symbol(),
		w.DerivePrivateKey(),
		w.DerivePublicKey(),
		w.DeriveAddress())
}
