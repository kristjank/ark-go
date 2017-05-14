package cmd

import (
	"ark-go/arkcoin"
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

//GetKey returns public/private key pair
func GetKey(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	b := h.Sum(nil)

	privKey, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), b)

	fmt.Println("Private Key :")
	fmt.Printf("%x \n", privKey.Serialize())
	fmt.Println("Public Key :")
	fmt.Printf("%x \n", pubKey.SerializeCompressed())
	//fmt.Printf(privKey.X)

	//GetAddress(pubKey.SerializeCompressed(), byte('\x17'))

	return ""
}

func GetAddress(password string) string {
	h := sha256.New()
	h.Write([]byte("this is a top secret passphrase"))
	b := h.Sum(nil)

	key1 := arkcoin.NewPrivateKey(b, arkcoin.ArkCoinMain)
	//arkcoin.NewPrivateKey(pb, param)
	key := arkcoin.NewPrivateKeyFromPassword(password, arkcoin.ArkCoinMain)

	//adr := key.PublicKey.Address()

	fmt.Printf("%x \n", key1.PublicKey.Serialize())
	//fmt.Printf(key.PublicKey.Address())

	return key.PublicKey.Address()
}
