package core

import (
	"ark-go/arkcoin"
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
		t.Error("Amount wrong")
	}

	if tx.Timestamp == 1 {
		if tx.Signature != "30450221008b7bc816d2224e34de8dac3dbe7d17789cf74f088a442a38f6e20fac632675bb02202d13119c896a2e282504341870d59cffe431395242834cd4d36afb62fbe27f97" {
			t.Error("Wrong signature")
		}

		if tx.SenderPublicKey != "034151a3ec46b5670a682b0a63394f863587d1bc97483b1b6c70eb58e7f0aed192" {
			t.Error("Wrong Public Key")
		}

		if tx.ID != "ccff05469c35db9091dcfb2fdb02b14dbf1b699f95a1ef4123ab891921e4b876" {
			t.Error("Wrong TX  ID")
		}
	}
	log.Println(t.Name(), "Transaction created OK, Json: ", tx.ToJSON())
}

func TestVerifyTransaction(t *testing.T) {
	tx := CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "")

	err := tx.Verify()
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(t.Name(), "Success")
}

func TestSecondVerifyTransaction(t *testing.T) {
	tx := CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "second top secret")

	err := tx.SecondVerify()
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(t.Name(), "Success")
}

func TestPostTransaction(t *testing.T) {
	/*tx := CreateTransaction("AUgTuukcKeE4XFdzaK6rEHMD5FLmVBSmHk",
		1,
		"ARK-GOLang is saying whoop whooop",
		"ski rose knock live elder parade dose device fetch betray loan holiday", "")

	arkapi := NewArkClient(nil)

	res, httpresponse, err := arkapi.PostTransaction(tx)
	if res.Success {
		log.Println(t.Name(), "Success,", httpresponse.Status, res.TransactionIDs)

	} else {
		log.Println(res.Message, res.Error, httpresponse.Status)
		t.Error(err.Error(), res.Error)
	}*/
}

func TestListTransaction(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := TransactionQueryParams{Limit: 10, SenderID: "AQLUKKKyKq5wZX7rCh4HJ4YFQ8bpTpPJgK"}

	transResponse, _, err := arkapi.ListTransaction(params)
	if transResponse.Success {
		log.Println(t.Name(), "Success, returned", transResponse.Count, "transactions")
	} else {
		t.Error(err.Error())
	}
}

func TestListTransactionUncomfirmed(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := TransactionQueryParams{Limit: 10, SenderID: "AQLUKKKyKq5wZX7rCh4HJ4YFQ8bpTpPJgK"}

	transResponse, _, err := arkapi.ListTransactionUnconfirmed(params)
	if transResponse.Success {
		log.Println(t.Name(), "Success, returned", transResponse.Count, "transactions")
	} else {
		t.Error(err.Error())
	}
}

func TestGetTransaction(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := TransactionQueryParams{ID: "bb032f1063fdd60844c250d3b76adcef3a75e686a0db2ef61be7e77ea0b8d293"}

	transResponse, _, err := arkapi.GetTransaction(params)
	if transResponse.Success {
		log.Println(t.Name(), "Success, returned tx with desc: ", transResponse.SingleTransaction.VendorField, "transactions")
	} else {
		log.Println(err.Error(), transResponse.Error)
		t.Error(err.Error())
	}
}

func TestGetTransactionUnconfirmed(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := TransactionQueryParams{SenderID: "AQLUKKKyKq5wZX7rCh4HJ4YFQ8bpTpPJgK", ID: "2105869df411b4fffd14eaf3bae10715acd176e7ea4a41df4141b35e717f2d39"}

	transResponse, _, err := arkapi.GetTransactionUnconfirmed(params)
	if transResponse.Success {
		log.Println(t.Name(), "Success, returned tx with desc: ", transResponse.SingleTransaction.VendorField, "transactions")
	} else {
		log.Println(err.Error(), transResponse.Error)
		t.Error(err.Error())
	}
}

func TestCreateDelegate(t *testing.T) {
	tx := CreateDelegate("chris", "this is a top secret passphrase", "")

	err := tx.Verify()
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(t.Name(), "Success")

}

func TestCreateVote(t *testing.T) {
	tx := CreateVote("+", "034151a3ec46b5670a682b0a63394f863587d1bc97483b1b6c70eb58e7f0aed192", "this is a top secret passphrase", "")

	err := tx.Verify()
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(t.Name(), "Success")

}

func TestCreateSecondSignature(t *testing.T) {
	tx := CreateSecondSignature("this is a top secret passphrase", "this is new second passphrase")

	err := tx.Verify()
	if err != nil {
		t.Error(err.Error())
	}
	log.Println(t.Name(), "Success")

}

func TestAddress(t *testing.T) {
	key := arkcoin.NewPrivateKeyFromPassword("this is a top secret passphrase", arkcoin.ActiveCoinConfig)

	if EnvironmentParams.Network.Type == MAINNET {
		if key.PublicKey.Address() != "AGeYmgbg2LgGxRW2vNNJvQ88PknEJsYizC" {
			t.Error("Address generation failed. Generated Address: ", key.PublicKey.Address())
		}
	}
	if EnvironmentParams.Network.Type == DEVNET {
		if key.PublicKey.Address() != "D61mfSggzbvQgTUe6JhYKH2doHaqJ3Dyib" {
			t.Error("Address generation failed. Generated Address: ", key.PublicKey.Address())
		}
	}
}
