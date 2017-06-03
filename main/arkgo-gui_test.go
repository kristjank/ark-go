package main

import (
	"ark-go/arkcoin"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"log"
	"testing"
)

func TestReadAccountData(t *testing.T) {
	pass := "this is key test password"
	b := make([]byte, 32)
	rand.Read(b)

	key := arkcoin.NewPrivateKeyFromPassword(pass, arkcoin.ArkCoinMain)

	ciphertext, err := encrypt([]byte(key.WIFAddress()), b)
	if err != nil {
		log.Fatal(err)
	}

	plaintext, err := decrypt(ciphertext, b)
	if err != nil {
		log.Fatal(err)
	}

	key1, err := arkcoin.FromWIF(string(plaintext), arkcoin.ArkCoinMain)
	if err != nil {
		log.Println(t.Name(), err.Error())
	}

	log.Println(key.PublicKey.Address(), key.PrivateKey.Serialize())
	log.Println(key1.PublicKey.Address(), key1.PrivateKey.Serialize())

	if key1.PublicKey.Address() != key.PublicKey.Address() {
		t.Error("Keys dont match")
	}
	//fmt.Printf("%x => %s\n", ciphertext, plaintext)
}

func TestGetSystemEnv(t *testing.T) {
	a := getSystemEnv()
	trHashBytes := sha256.New()
	trHashBytes.Write([]byte(a))
	log.Println(hex.EncodeToString(trHashBytes.Sum(nil)))
}

func TestSave(t *testing.T) {
	pass := "password"
	key := arkcoin.NewPrivateKeyFromPassword(pass, arkcoin.ArkCoinMain)
	log.Println(key.PublicKey.Address(), key.PrivateKey.Serialize())

	save(pass)

	key1, _ := read()
	log.Println(key1.PublicKey.Address(), key1.PrivateKey.Serialize())

	if key1.PublicKey.Address() != key.PublicKey.Address() {
		t.Error("Keys don't match")
	}
}
