package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/kristjank/ark-go/arkcoin"
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
	/*pass := "password"

	save(pass, "")

	p1, _ := read()
	if p1 != pass {
		t.Error("Keys don't match")
	}*/
}

func TestSave1(t *testing.T) {
	/*pass := "password"

	save(pass, pass)

	p1, p2 := read()

	if p1 != pass {
		t.Error("Keys1 don't match")
	}

	if p2 != pass {
		t.Error("Keys2 don't match")
	}*/
}

func TestCreateLogFolder(t *testing.T) {
	tt := time.Now()

	folderName := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d",
		tt.Year(), tt.Month(), tt.Day(),
		tt.Hour(), tt.Minute(), tt.Second())
	log.Println("log/" + folderName)

	err := os.MkdirAll("log/"+folderName, os.ModePerm)
	if err != nil {
		t.Error(err.Error())
	}
}
