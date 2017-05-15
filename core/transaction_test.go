package core

import (
	"encoding/json"
	"log"
	"testing"
)

func TestCreateTransaction(t *testing.T) {

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

	txJSON, err := json.Marshal(tx)
	if err == nil {
		log.Println(txJSON)
	}
}

func TestSignTransaction(t *testing.T) {
	tx := CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "")

	if tx.Amount == 0 {
		t.Error("Amount  is zero")
	}

	log.Println("Tx PubKey: ", tx.SenderPublicKey)
	log.Println("Tx Signature: ", tx.Signature)
}

/*func TestToBytes(*testing.T) {
	tx := CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "")

	b := tx.ToBytes(true, true)
	if len(b) == 0 {
		log.Fatal("Error")
	}
	encodedStr := hex.EncodeToString(b)
	//log.Printf("%x \n", b)
	log.Println("ToBytes: ", encodedStr)
}*/
