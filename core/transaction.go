package core

import (
	"ark-go/arkcoin"
	"ark-go/arkcoin/base58"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"log"
)

//Transaction struct - represents structure of ARK.io blockchain transaction
type Transaction struct {
	Timestamp          int
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

func HexDecodeData(data string) []byte {
	src := []byte(data)

	dst := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	return dst[:n]
}

//ToBytes returns bytearray of the Transaction object to be signed and send to blockchain
func ToBytes(tx *Transaction, skipSignature, skipSecondSignature bool) []byte {
	txBuf := new(bytes.Buffer)
	binary.Write(txBuf, binary.LittleEndian, tx.Type)
	binary.Write(txBuf, binary.LittleEndian, tx.Timestamp)
	binary.Write(txBuf, binary.LittleEndian, HexDecodeData(tx.SenderPublicKey))

	if tx.RequesterPublicKey != "" {
		res, err := base58.Decode(tx.RequesterPublicKey)
		if err != nil {
			binary.Write(txBuf, binary.LittleEndian, res)
		}
	}

	if tx.RecipientID != "" {
		res, err := base58.Decode(tx.RecipientID)
		if err != nil {
			binary.Write(txBuf, binary.LittleEndian, res)
		}
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

	binary.Write(txBuf, binary.LittleEndian, tx.Amount)
	binary.Write(txBuf, binary.LittleEndian, tx.Fee)

	switch tx.Type {
	case 1:
		binary.Write(txBuf, binary.LittleEndian, HexDecodeData(tx.Signature))
	case 2:
		//buffer.Put(Encoding.ASCII.GetBytes(asset["username"]));
	case 3:
		//votes
	}

	if !skipSignature && len(tx.Signature) > 0 {
		binary.Write(txBuf, binary.LittleEndian, HexDecodeData(tx.Signature))
	}

	if !skipSecondSignature && len(tx.SignSignature) > 0 {
		binary.Write(txBuf, binary.LittleEndian, HexDecodeData(tx.SignSignature))
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
	Sign(&tx, passphrase)

	//if (secondPassphrase != nill)
	//	tx.SecondSign(secondPassphrase);

	//tx.Id = Crypto.GetId(tx);
	return &tx
}

//Sign the Transaction
func Sign(tx *Transaction, passphrase string) {
	key := arkcoin.NewPrivateKeyFromPassword(passphrase, arkcoin.ArkCoinMain)

	tx.SenderPublicKey = hex.EncodeToString(key.PublicKey.Serialize())

	trHashBytes := sha256.New()
	trHashBytes.Write(ToBytes(tx, true, true))
	trHashBytes.Sum(nil)

	sig, err := key.Sign(trHashBytes.Sum(nil))
	if err == nil {
		tx.Signature = hex.EncodeToString(sig)
	}
}
