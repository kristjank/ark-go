package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"

	"github.com/kristjank/ark-go/arkcoin"
	"github.com/kristjank/ark-go/arkcoin/base58"
)

type TransactionType byte

const (
	SENDARK         = 0
	SECONDSIGNATURE = 1
	CREATEDELEGATE  = 2
	VOTE            = 3
	MULTISIGNATURE  = 4
)

//Transaction struct - represents structure of ARK.io blockchain transaction
//It is used to post transaction to mainnet and to receive results from arkapi
//Empty fields are emmited by default
type Transaction struct {
	ID                    string            `json:"id,omitempty"`
	Timestamp             int32             `json:"timestamp,omitempty"`
	RecipientID           string            `json:"recipientId,omitempty"`
	Amount                int64             `json:"amount,omitempty"`
	Asset                 map[string]string `json:"asset,omitempty"`
	Fee                   int64             `json:"fee,omitempty"`
	Type                  byte              `json:"type"`
	VendorField           string            `json:"vendorField,omitempty"`
	Signature             string            `json:"signature,omitempty"`
	SignSignature         string            `json:"signSignature,omitempty"`
	SenderPublicKey       string            `json:"senderPublicKey,omitempty"`
	SecondSenderPublicKey string            `json:"secondSenderPublicKey,omitempty"`
	RequesterPublicKey    string            `json:"requesterPublicKey,omitempty"`
	Blockid               string            `json:"blockid,omitempty"`
	Height                int               `json:"height,omitempty"`
	SenderID              string            `json:"senderId,omitempty"`
	Confirmations         int               `json:"confirmations,omitempty"`
}

func fromBytes(txbytes []byte) Transaction {
	txReader := bytes.NewReader(txbytes)

	tx := Transaction{}

	//Get Transaction Type
	binary.Read(txReader, binary.LittleEndian, &tx.Type)

	//GetTimeStamp from 	binary.Write(txBuf, binary.LittleEndian, uint32(tx.Timestamp))
	binary.Read(txReader, binary.LittleEndian, &tx.Timestamp)
	//	tx.Timestamp = timestamp
	//txReader.
	return tx
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
	case SECONDSIGNATURE:
		binary.Write(txBuf, binary.LittleEndian, quickHexDecode(tx.Asset["signature"]))
	case CREATEDELEGATE:
		usernameBytes := []byte(tx.Asset["username"])
		binary.Write(txBuf, binary.LittleEndian, usernameBytes)
	case VOTE:
		voteBytes := []byte(tx.Asset["votes"])
		binary.Write(txBuf, binary.LittleEndian, voteBytes)
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
	tx := Transaction{
		Type:        SENDARK,
		RecipientID: recipientID,
		Amount:      satoshiAmount,
		Fee:         EnvironmentParams.Fees.Send,
		VendorField: vendorField,
	}

	tx.Timestamp = GetTime() //1
	tx.sign(passphrase)

	if len(secondPassphrase) > 0 {
		tx.secondSign(secondPassphrase)
	}

	tx.getID() //calculates id of transaction
	return &tx
}

//CreateVote transaction used to vote for a chosen Delegate
//if updown value = "+" vot is given to the specified PublicKey
//if updown value = "-" vot is taken from the specified PublicKey
func CreateVote(updown, delegatePubKey, passphrase, secondPassphrase string) *Transaction {
	tx := Transaction{
		Type:        VOTE,
		Fee:         EnvironmentParams.Fees.Vote,
		VendorField: "Delegate vote transaction",
		Asset:       make(map[string]string),
	}
	key := arkcoin.NewPrivateKeyFromPassword(passphrase, arkcoin.ActiveCoinConfig)
	tx.RecipientID = key.PublicKey.Address()

	tx.Asset["votes"] = updown + delegatePubKey
	tx.Timestamp = GetTime() //1
	tx.sign(passphrase)

	if len(secondPassphrase) > 0 {
		tx.secondSign(secondPassphrase)
	}

	tx.getID() //calculates id of transaction
	return &tx
}

//CreateDelegate creates and returns new Transaction struct...
func CreateDelegate(username, passphrase, secondPassphrase string) *Transaction {
	tx := Transaction{
		Type:        CREATEDELEGATE,
		Fee:         EnvironmentParams.Fees.Delegate,
		VendorField: "Create delegate tx",
		Asset:       make(map[string]string),
	}
	tx.Asset["username"] = username
	tx.Timestamp = GetTime() //1
	tx.sign(passphrase)

	if len(secondPassphrase) > 0 {
		tx.secondSign(secondPassphrase)
	}

	tx.getID() //calculates id of transaction
	return &tx
}

//CreateSecondSignature creates and returns new Transaction struct...
func CreateSecondSignature(passphrase, secondPassphrase string) *Transaction {
	tx := Transaction{
		Type:        SECONDSIGNATURE,
		Fee:         EnvironmentParams.Fees.SecondSignature,
		VendorField: "Create second signature",
		Asset:       make(map[string]string),
	}

	key := arkcoin.NewPrivateKeyFromPassword(secondPassphrase, arkcoin.ActiveCoinConfig)
	tx.Asset["signature"] = hex.EncodeToString(key.PublicKey.Serialize())
	tx.Timestamp = GetTime() //1
	tx.sign(passphrase)

	tx.getID() //calculates id of transaction
	return &tx
}

//Sign the Transaction
func (tx *Transaction) sign(passphrase string) {
	key := arkcoin.NewPrivateKeyFromPassword(passphrase, arkcoin.ActiveCoinConfig)

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
	key := arkcoin.NewPrivateKeyFromPassword(passphrase, arkcoin.ActiveCoinConfig)

	tx.SecondSenderPublicKey = hex.EncodeToString(key.PublicKey.Serialize())
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
//if return == nill verification was succesfull
func (tx *Transaction) Verify() error {
	key, err := arkcoin.NewPublicKey(quickHexDecode(tx.SenderPublicKey), arkcoin.ActiveCoinConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	trHashBytes := sha256.New()
	trHashBytes.Write(tx.toBytes(true, true))
	return key.Verify(quickHexDecode(tx.Signature), trHashBytes.Sum(nil))

}

//SecondVerify function verifies if tx is validly signed
//if return == nill verification was succesfull
func (tx *Transaction) SecondVerify() error {
	key, err := arkcoin.NewPublicKey(quickHexDecode(tx.SecondSenderPublicKey), arkcoin.ActiveCoinConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	trHashBytes := sha256.New()
	trHashBytes.Write(tx.toBytes(false, true))
	return key.Verify(quickHexDecode(tx.SignSignature), trHashBytes.Sum(nil))
}

//PostTransactionResponse structure for call /peer/list
type PostTransactionResponse struct {
	Success        bool     `json:"success"`
	Message        string   `json:"message"`
	Error          string   `json:"error"`
	TransactionIDs []string `json:"transactionIds"`
}

//TransactionPayload - list of tx to send to network
type TransactionPayload struct {
	Transactions []*Transaction `json:"transactions"`
}

//TransactionQueryParams for returing filtered list of transactions
type TransactionQueryParams struct {
	ID          string          `url:"id,omitempty"`
	BlockID     string          `url:"blockId,omitempty"`
	SenderID    string          `url:"senderId,omitempty"`
	RecipientID string          `url:"recipientId,omitempty"`
	Limit       int             `url:"limit,omitempty"`
	Offset      int             `url:"offset,omitempty"`
	OrderBy     string          `url:"orderBy,omitempty"` //"Name of column to order. After column name must go 'desc' or 'asc' to choose order type, prefix for column name is t_. Example: orderBy=t_timestamp:desc (String)"
	Type        TransactionType `url:"type,omitempty"`
}

//TransactionResponse structure holds parsed jsong reply from ark-node
//when calling list methods the Transactions [] has results
//when calling get methods the transaction object (Single) has results
type TransactionResponse struct {
	Success           bool          `json:"success"`
	Transactions      []Transaction `json:"transactions"`
	SingleTransaction Transaction   `json:"transaction"`
	Count             string        `json:"count"`
	Error             string        `json:"error"`
}

//PostTransaction to selected ARKNetwork
func (s *ArkClient) PostTransaction(payload TransactionPayload) (PostTransactionResponse, *http.Response, error) {
	respTr := new(PostTransactionResponse)
	errTr := new(ArkApiResponseError)

	/*var payload transactionPayload
	payload.Transactions = append(payload.Transactions, tx)
	*/
	resp, err := s.sling.New().Post("peer/transactions").BodyJSON(payload).Receive(respTr, errTr)

	if err == nil {
		err = errTr
	}

	return *respTr, resp, err
}

//ListTransaction function returns list of peers from ArkNode
func (s *ArkClient) ListTransaction(params TransactionQueryParams) (TransactionResponse, *http.Response, error) {
	transactionResponse := new(TransactionResponse)
	transactionResponseErr := new(ArkApiResponseError)
	resp, err := s.sling.New().Get("api/transactions").QueryStruct(&params).Receive(transactionResponse, transactionResponseErr)
	if err == nil {
		err = transactionResponseErr
	}

	return *transactionResponse, resp, err
}

//ListTransactionUnconfirmed function returns list of peers from ArkNode
func (s *ArkClient) ListTransactionUnconfirmed(params TransactionQueryParams) (TransactionResponse, *http.Response, error) {
	transactionResponse := new(TransactionResponse)
	transactionResponseErr := new(ArkApiResponseError)
	resp, err := s.sling.New().Get("api/transactions/unconfirmed").QueryStruct(&params).Receive(transactionResponse, transactionResponseErr)
	if err == nil {
		err = transactionResponseErr
	}

	return *transactionResponse, resp, err
}

//GetTransaction function returns list of peers from ArkNode
func (s *ArkClient) GetTransaction(params TransactionQueryParams) (TransactionResponse, *http.Response, error) {
	transactionResponse := new(TransactionResponse)
	transactionResponseErr := new(ArkApiResponseError)
	resp, err := s.sling.New().Get("api/transactions/get").QueryStruct(&params).Receive(transactionResponse, transactionResponseErr)
	if err == nil {
		err = transactionResponseErr
	}

	return *transactionResponse, resp, err
}

//GetTransactionUnconfirmed function returns list of peers from ArkNode
func (s *ArkClient) GetTransactionUnconfirmed(params TransactionQueryParams) (TransactionResponse, *http.Response, error) {
	transactionResponse := new(TransactionResponse)
	transactionResponseErr := new(ArkApiResponseError)
	resp, err := s.sling.New().Get("api/transactions/unconfirmed/get").QueryStruct(&params).Receive(transactionResponse, transactionResponseErr)
	if err == nil {
		err = transactionResponseErr
	}

	return *transactionResponse, resp, err
}
