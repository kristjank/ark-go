package arkcoin

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

//ToBytes returns bytearray of the Transaction object to be signed and send to blockchain
func ToBytes(skipSignature, skipSecondSignature bool) []byte {
	return nil
}

//CreateTransaction creates and returns new Transaction struct...
func CreateTransaction(recipientID string, satoshiAmount int64, vendorField, passphrase, secondPassphrase string) Transaction {
	tx := Transaction{Type: 0, RecipientID: recipientID, Amount: satoshiAmount, Fee: 10000000, VendorField: vendorField}

	tx.Timestamp = 1 //Slot.GetTime();
	//tx.Sign(passphrase);
	//tx.StrBytes = Encoders.Hex.EncodeData(tx.ToBytes());
	//if (secondPassphrase != nill)
	//	tx.SecondSign(secondPassphrase);

	//tx.Id = Crypto.GetId(tx);
	return tx
}
