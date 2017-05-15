package core

import (
	"log"
	"testing"
)

func TestCreateSignTransaction(t *testing.T) {
	tx := CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "")

	if tx.Amount == 0 {
		t.Error("Amount is zero")
	}

	if tx.Amount != 133380000000 {
		log.Fatal("Amount wrong")
	}

	if tx.Signature != "30450221008b7bc816d2224e34de8dac3dbe7d17789cf74f088a442a38f6e20fac632675bb02202d13119c896a2e282504341870d59cffe431395242834cd4d36afb62fbe27f97" {
		log.Fatal("Wrong signature")
	}

	if tx.SenderPublicKey != "034151a3ec46b5670a682b0a63394f863587d1bc97483b1b6c70eb58e7f0aed192" {
		log.Fatal("Wrong Public Key")
	}

	if tx.ID != "ccff05469c35db9091dcfb2fdb02b14dbf1b699f95a1ef4123ab891921e4b876" {
		log.Fatal("Wrong TX  ID")
	}

	log.Println(tx.ToJSON())
}
