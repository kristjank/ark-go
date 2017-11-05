package main

import (
	"time"

	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func runTests() {
	testRecord := createTestRecord()
	testRecord.Save()

	testRecord.TestStarted = time.Now()
	for xx := 0; xx < viper.GetInt("env.txIterations"); xx++ {
		ArkAPIClient = ArkAPIClient.SetActiveConfiguration(core.DEVNET)
		payload := core.TransactionPayload{}

		for i := 0; i < viper.GetInt("env.txPerPayload"); i++ {
			tx := core.CreateTransaction(viper.GetString("account.recepient"),
				int64(i+1),
				viper.GetString("env.txDescription"),
				viper.GetString("account.passphrase"), viper.GetString("account.secondPassphrase"))
			payload.Transactions = append(payload.Transactions, tx)
		}

		log.Info("Sending transactions to ", ArkAPIClient.GetActivePeer(), "nr of tx: ", len(payload.Transactions))

		testIterRecord := createTestIterationRecord(testRecord.ID)
		testIterRecord.Save()
		res, httpresponse, err := ArkAPIClient.PostTransaction(payload)
		testIterRecord.IterationStopped = time.Now()

		if res.Success {
			log.Info("Success,", httpresponse.Status, xx)
			testIterRecord.TestStatus = "SUCCESS"
			testIterRecord.TxIDs = res.TransactionIDs
		} else {
			testIterRecord.TestStatus = "FAILED"
			if httpresponse != nil {
				log.Error(res.Message, res.Error, xx)
			}
			log.Error(err.Error(), res.Error)
		}
		testIterRecord.Update()
	}
	testRecord.TestStopped = time.Now()
	testRecord.Update()
	log.Info("The call took %v to run.\n", testRecord.TestStopped.Sub(testRecord.TestStarted))

}
