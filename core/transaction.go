package core

import (
	"ark-go/arkcoin"
	"ark-go/arkcoin/base58"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
)

//Transaction struct - represents structure of ARK.io blockchain transaction
type Transaction struct {
	Timestamp          int32
	RecipientID        string
	Amount             int64
	Fee                int64
	Type               byte
	VendorField        string
	Signature          string
	SignSignature      string
	SenderPublicKey    string
	RequesterPublicKey string
	ID                 string
}

//ToBytes returns bytearray of the Transaction object to be signed and send to blockchain
func (tx *Transaction) toBytes(skipSignature, skipSecondSignature bool) []byte {
	txBuf := new(bytes.Buffer)
	binary.Write(txBuf, binary.LittleEndian, tx.Type)
	binary.Write(txBuf, binary.LittleEndian, uint32(tx.Timestamp))

	binary.Write(txBuf, binary.LittleEndian, quickHexDecode(tx.SenderPublicKey))

	if tx.RequesterPublicKey != "" {
		res, err := base58.Decode(tx.RequesterPublicKey)
		if err != nil {
			binary.Write(txBuf, binary.LittleEndian, res)
		}
	}

	if tx.RecipientID != "" {
		res, err := base58.Decode(tx.RecipientID)
		if err != nil {
			log.Fatal("Error converting Decoding b58 ", err.Error())
		}
		binary.Write(txBuf, binary.LittleEndian, res)
	} else {
		binary.Write(txBuf, binary.LittleEndian, make([]byte, 21))
	}

	if tx.VendorField != "" {
		vendorBytes := []byte(tx.VendorField)
		if len(vendorBytes) < 65 {
			binary.Write(txBuf, binary.LittleEndian, vendorBytes)

			bs := make([]byte, 64-len(vendorBytes))
			binary.Write(txBuf, binary.LittleEndian, bs)
		}
	} else {
		binary.Write(txBuf, binary.LittleEndian, make([]byte, 64))
	}

	binary.Write(txBuf, binary.LittleEndian, uint64(tx.Amount))
	binary.Write(txBuf, binary.LittleEndian, uint64(tx.Fee))

	switch tx.Type {
	case 1:
		binary.Write(txBuf, binary.LittleEndian, quickHexDecode(tx.Signature))
	case 2:
		//buffer.Put(Encoding.ASCII.GetBytes(asset["username"]));
	case 3:
		//votes
	}

	if !skipSignature && len(tx.Signature) > 0 {
		binary.Write(txBuf, binary.LittleEndian, quickHexDecode(tx.Signature))
	}

	if !skipSecondSignature && len(tx.SignSignature) > 0 {
		binary.Write(txBuf, binary.LittleEndian, quickHexDecode(tx.SignSignature))
	}

	return txBuf.Bytes()
}

//CreateTransaction creates and returns new Transaction struct...
func CreateTransaction(recipientID string, satoshiAmount int64, vendorField, passphrase, secondPassphrase string) *Transaction {
	tx := Transaction{Type: 0,
		RecipientID: recipientID,
		Amount:      satoshiAmount,
		Fee:         arkcoin.ArkCoinMain.Fees.Send,
		VendorField: vendorField}

	tx.Timestamp = 1 //Slot.GetTime();
	tx.sign(passphrase)

	if len(secondPassphrase) > 0 {
		tx.secondSign(secondPassphrase)
	}

	tx.getID() //calculates id of transaction
	return &tx
}

//Sign the Transaction
func (tx *Transaction) sign(passphrase string) {
	key := arkcoin.NewPrivateKeyFromPassword(passphrase, arkcoin.ArkCoinMain)

	tx.SenderPublicKey = hex.EncodeToString(key.PublicKey.Serialize())

	trHashBytes := sha256.New()
	trHashBytes.Write(tx.toBytes(true, true))

	sig, err := key.Sign(trHashBytes.Sum(nil))
	if err == nil {
		tx.Signature = hex.EncodeToString(sig)
	}
}

//SecondSign the Transaction
func (tx *Transaction) secondSign(passphrase string) {
	key := arkcoin.NewPrivateKeyFromPassword(passphrase, arkcoin.ArkCoinMain)

	trHashBytes := sha256.New()
	trHashBytes.Write(tx.toBytes(false, true))

	sig, err := key.Sign(trHashBytes.Sum(nil))
	if err == nil {
		tx.SignSignature = hex.EncodeToString(sig)
	}
}

//GetID returns calculated ID of trancation - hashed s256
func (tx *Transaction) getID() {
	trHashBytes := sha256.New()
	trHashBytes.Write(tx.toBytes(false, false))

	tx.ID = hex.EncodeToString(trHashBytes.Sum(nil))
}

//ToJSON converts transaction object to JSON string
func (tx *Transaction) ToJSON() string {
	txJSON, err := json.Marshal(tx)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(txJSON)
}

func quickHexDecode(data string) []byte {
	res, err := hex.DecodeString(data)
	if err != nil {
		log.Fatal(err.Error())
	}
	return res
}

//Verify function verifies if tx is validly signed
func (tx *Transaction) Verify() bool {
	return tx.verifyHelper(true)
}

//SecondVerify function verifies if tx is validly signed
func (tx *Transaction) SecondVerify() bool {
	return tx.verifyHelper(false)
}

func (tx *Transaction) verifyHelper(first bool) bool {
	key, err := arkcoin.NewPublicKey(quickHexDecode(tx.SenderPublicKey), arkcoin.ArkCoinMain)
	if err != nil {
		log.Fatal(err.Error())
	}

	trHashBytes := sha256.New()
	trHashBytes.Write(tx.toBytes(first, true))

	var res
	if first {
		res := key.Verify(quickHexDecode(tx.Signature), trHashBytes.Sum(nil))
	} else {
		res := key.Verify(quickHexDecode(tx.SignSignature), trHashBytes.Sum(nil))
	}

	if res == nil {
		return true
	}
	return false

}
