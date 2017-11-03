package main

import (
	"log"
	"time"

	"github.com/kristjank/ark-go/core"
)

func main() {
	arkapi := core.NewArkClient(nil)
	arkapi = arkapi.SetActiveConfiguration(core.DEVNET)

	recepient := "AUgTuukcKeE4XFdzaK6rEHMD5FLmVBSmHk"
	passphrase := "ski rose knock live elder parade dose device fetch betray loan holiday"

	if core.EnvironmentParams.Network.Type == core.DEVNET {
		recepient = "DFTzLwEHKKn3VGce6vZSueEmoPWpEZswhB"
		passphrase = "outer behind tray slice trash cave table divert wild buddy snap news"
	}
	t0 := time.Now()

	for xx := 0; xx < 250; xx++ {
		arkapi = arkapi.SetActiveConfiguration(core.DEVNET)

		var payload core.TransactionPayload

		for i := 0; i < 300; i++ {
			tx := core.CreateTransaction(recepient,
				int64(i+1),
				"1ARK-GOLang is saying whoop whooop",
				passphrase, "")
			payload.Transactions = append(payload.Transactions, tx)
		}

		res, httpresponse, err := arkapi.PostTransaction(payload)
		if res.Success {
			log.Println("Success,", httpresponse.Status, xx)

		} else {
			if httpresponse != nil {
				log.Println(res.Message, res.Error, xx)
			}
			log.Println(err.Error(), res.Error)
		}
		payload.Transactions = nil
		time.Sleep(1000)
	}

	t1 := time.Now()
	log.Printf("The call took %v to run.\n", t1.Sub(t0))
}
