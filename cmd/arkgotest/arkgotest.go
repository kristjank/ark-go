package main

import (
	"io"
	"os"
	"time"

	"github.com/kristjank/ark-go/core"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func init() {
	initLogger()
	loadConfig()

}
func initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("log/arkgotest.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(file))
	} else {
		log.Error("Failed to log to file, using default stderr")
	}

}

func loadConfig() {
	viper.AddConfigPath("cfg")  // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		log.Info("No productive config found - loading sample")
		// try to load sample config
		viper.SetConfigName("sample.config")
		viper.AddConfigPath("cfg")
		err := viper.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			log.Fatal("No configuration file loaded - using defaults")
		}
	}

	viper.SetDefault("env.txPerPayload", 250)
	viper.SetDefault("env.txMultiBroadCast", 1)
	viper.SetDefault("env.dbfilename", "db/testlog.db")
}

func checkTx4Confirmations() {

}

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

	for xx := 0; xx < 1; xx++ {
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
