package arkcoin

import "testing"

func TestToBytes(t *testing.T) {
	bytes := ToBytes(false, false)
	if bytes != nil {
		t.Error("Bytes not nill")
	}
}

func TestCreateTransaction(t *testing.T) {
	tx := CreateTransaction("AXoXnFi4z1Z6aFvjEYkDVCtBGW2PaRiM25",
		133380000000,
		"This is first transaction from ARK-NET",
		"this is a top secret passphrase", "")

	if tx.Amount == 0 {
		t.Error("Amount is zero")
	}
}
