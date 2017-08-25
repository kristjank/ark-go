package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func floatEquals(a, b float64) bool {
	EPSILON := 0.000000000000001
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func checkConfigSharingRatio() bool {
	a1 := viper.GetFloat64("voters.shareratio")
	a2 := viper.GetFloat64("costs.shareratio")
	a3 := viper.GetFloat64("reserve.shareratio")
	a4 := viper.GetFloat64("personal.shareratio")

	if !floatEquals(a1+a2+a3+a4, 1.0) {
		log.Info("Wrong config. Check share ration percentages!")
		return false
	}
	return true
}

func log2csv(payload core.TransactionPayload, txids []string, filecsv *os.File, status string) {
	records := [][]string{
		{"ADDRESS", "SENT AMOUNT", "TimeStamp", "TxId", "ApiResponse"},
	}

	for ix, el := range payload.Transactions {
		//		sAmount := fmt.Sprintf("%15.8f", float64(el.Amount)/float64(core.SATOSHI))
		timeTx := core.GetTransactionTime(el.Timestamp)
		localTime := timeTx.Local()

		var line []string
		if txids != nil {
			line = []string{el.RecipientID, strconv.FormatFloat(float64(el.Amount)/float64(core.SATOSHI), 'f', -1, 64), localTime.Format("2006-01-02 15:04:05"), txids[ix], status}
		} else {
			line = []string{el.RecipientID, strconv.FormatFloat(float64(el.Amount)/float64(core.SATOSHI), 'f', -1, 64), localTime.Format("2006-01-02 15:04:05"), "N/A", status}
		}

		records = append(records, line)

	}
	w := csv.NewWriter(filecsv)
	defer w.Flush()
	w.WriteAll(records)
}

func getSystemEnv() string {
	var buffer bytes.Buffer
	buffer.WriteString(os.Getenv("OS"))
	buffer.WriteString(os.Getenv("PROCESSOR_ARCHITECTURE"))
	buffer.WriteString(os.Getenv("PROCESSOR_IDENTIFIER"))
	buffer.WriteString(os.Getenv("COMPUTERNAME"))
	buffer.WriteString(os.Getenv("ComSpec"))

	buffer.WriteString(os.Getenv("OS"))
	buffer.WriteString(os.Getenv("PROCESSOR_ARCHITECTURE"))
	buffer.WriteString(os.Getenv("PROCESSOR_IDENTIFIER"))
	buffer.WriteString(os.Getenv("COMPUTERNAME"))
	buffer.WriteString(os.Getenv("ComSpec"))

	return buffer.String()
}

func save(p1, p2 string) {
	ciphertext, _ := encrypt([]byte(p1), getRandHash())
	ioutil.WriteFile("assembly.ark", ciphertext, 0644)

	if p2 != "" {
		ciphertext, err := encrypt([]byte(p2), getRandHash())
		if err != nil {
			log.Info("Error encrypting")
		}
		ioutil.WriteFile("assembly1.ark", ciphertext, 0644)
	} else {
		os.Remove("assembly1.ark")
	}
}

/*func read() (*arkcoin.PrivateKey, *arkcoin.PrivateKey) {
	dat, err := ioutil.ReadFile("assembly.ark")
	if err != nil {
		log.Info(err.Error())
	}
	plaintext, _ := decrypt(dat, getRandHash())
	key1 := arkcoin.NewPrivateKeyFromPassword(string(plaintext), arkcoin.ActiveCoinConfig)

	var key2 *arkcoin.PrivateKey
	if _, err := os.Stat("assembly1.ark"); err == nil {
		dat, _ = ioutil.ReadFile("assembly1.ark")
		plaintext, _ = decrypt(dat, getRandHash())
		key2 = arkcoin.NewPrivateKeyFromPassword(string(plaintext), arkcoin.ActiveCoinConfig)
	}

	return key1, key2
}*/

func read() (string, string) {
	dat, err := ioutil.ReadFile("assembly.ark")
	if err != nil {
		log.Info(err.Error())
	}
	p1, _ := decrypt(dat, getRandHash())

	var p2 []byte

	if _, err := os.Stat("assembly1.ark"); err == nil {
		dat, _ = ioutil.ReadFile("assembly1.ark")
		p2, _ = decrypt(dat, getRandHash())
	}

	return string(p1), string(p2)
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func getRandHash() []byte {
	a := getSystemEnv()

	trHashBytes := sha256.New()
	trHashBytes.Write([]byte(a))

	return trHashBytes.Sum(nil)
}

func createLogFolder() string {
	tt := time.Now()

	folderName := fmt.Sprintf("%d-%02d-%02dT%02d-%02d-%02d",
		tt.Year(), tt.Month(), tt.Day(),
		tt.Hour(), tt.Minute(), tt.Second())
	log.Println("log/" + folderName)

	err := os.MkdirAll("log/"+folderName, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
	}

	return folderName
}
