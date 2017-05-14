package cmd

import (
	"ark-go/arkcoin"
	"ark-go/core"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

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

	fmt.Printf("%x \n", key1.Serialize())
	return key.PublicKey.Address()
}

func main() {
	tx := core.CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "")

	b, err := json.Marshal(&tx)
	if err != nil {
		log.Println(b)
	} else {
		log.Fatal(err.Error())
	}
}
