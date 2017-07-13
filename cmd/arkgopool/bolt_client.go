package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/kristjank/ark-go/core"
	"github.com/kristjank/goark-node/api/model"
	"github.com/spf13/viper"
)

//IBoltClient interface definition
type IBoltClient interface {
	OpenBoltDb()
	QueryPayments(recepientID string) error
	SavePayment(block model.Block) error
	InitializeBucket()
	GetAllPayments() (err error)
	Close()
}

var (
	ErrBucketNotFound = errors.New("Bucket not found")
	ErrKeyNotFound    = errors.New("Key not found")
	ErrDoesNotExist   = errors.New("Does not exist")
	ErrFoundIt        = errors.New("Found it")
	ErrExistsInSet    = errors.New("Element already exists in set")
	ErrInvalidID      = errors.New("Element ID can not contain \":\"")
)

//Constant names for BoltDb bucket initializations
const (
	PaymentBucket      = "PaymentBucket"
	VoterEarningBucket = "VoterEarningsBucket"
)

//BoltClient Realimplementation
type BoltClient struct {
	boltDB *bolt.DB
}

//OpenBoltDb db opening
func (bc *BoltClient) OpenBoltDb() {
	var err error
	bc.boltDB, err = bolt.Open(viper.GetString("dbFilename"), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Close the database
func (bc *BoltClient) Close() {
	bc.boltDB.Close()
}

//QueryPayment returns the block by id
func (bc *BoltClient) QueryPayment(transID string) (core.Transaction, error) {
	// Allocate an empty Account instance we'll let json.Unmarhal populate for us in a bit.
	tx := core.Transaction{}

	// Read an object from the bucket using boltDB.View
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		// Read the bucket from the DB
		b := tx.Bucket([]byte(PaymentBucket))

		// Read the value identified by our accountId supplied as []byte
		blockBytes := b.Get([]byte(transID))
		if blockBytes == nil {
			return fmt.Errorf("No block found for " + transID)
		}
		// Unmarshal the returned bytes into the account struct we created at
		// the top of the function
		json.Unmarshal(blockBytes, &tx)

		// Return nil to indicate nothing went wrong, e.g no error
		return nil
	})
	// If there were an error, return the error
	if err != nil {
		return core.Transaction{}, err
	}
	// Return the Account struct and nil as error.
	return tx, nil
}

//Check Naive healthcheck, just makes sure the DB connection has been initialized.
func (bc *BoltClient) Check() bool {
	return bc.boltDB != nil
}

//InitializeBucket Creates an "BlockBucket" in our BoltDB. It will overwrite any existing bucket of the same name.
func (bc *BoltClient) InitializeBucket() {
	bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(PaymentBucket))
		if err != nil {
			return fmt.Errorf("create Bucket failed: %s", err)
		}

		_, err = tx.CreateBucket([]byte(VoterEarningBucket))
		if err != nil {
			return fmt.Errorf("create Bucket failed: %s", err)
		}

		return nil
	})
}

//SaveTransaction to TransactionBucket
func (bc *BoltClient) SaveTransaction(trans core.Transaction) (string, error) {
	// Serialize the struct to JSON
	jsonBytes, _ := json.Marshal(trans)

	// Write the data to the BlockBucketBlockBucket
	err := bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PaymentBucket))
		err := b.Put([]byte(trans.ID), jsonBytes)
		return err
	})
	if err != nil {
		return "", err
	}
	return trans.ID, err

}

//GetAllPayments elements of a list
func (bc *BoltClient) GetAllPayments() (results []core.Transaction, err error) {
	return results, bc.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(PaymentBucket))
		if bucket == nil {
			return ErrBucketNotFound
		}
		return bucket.ForEach(func(_, value []byte) error {
			tx := core.Transaction{}
			json.Unmarshal(value, &tx)
			results = append(results, tx)
			return nil // Continue ForEach
		})
	})
}
