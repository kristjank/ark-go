package main

import (
	"time"

	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func runTests() {
	t0 := time.Now()
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
		res, httpresponse, err := ArkAPIClient.PostTransaction(payload)
		if res.Success {
			log.Info("Success,", httpresponse.Status, xx)

		} else {
			if httpresponse != nil {
				log.Error(res.Message, res.Error, xx)
			}
			log.Error(err.Error(), res.Error)
		}
		time.Sleep(1000)
	}

	t1 := time.Now()
	log.Info("The call took %v to run.\n", t1.Sub(t0))
}
