package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/asdine/storm"
	"github.com/kristjank/ark-go/core"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

//ArkAPIClient - ARKAPI Client
var ArkAPIClient *core.ArkClient

//ArkTestDB - testDB to store log records
var ArkTestDB *storm.DB

func init() {
	initLogger()
	loadConfig()
	ArkAPIClient = core.NewArkClient(nil)
	openDB()
}

func openDB() {
	log.Info("Opening/Reopening database")
	var err error
	ArkTestDB, err = storm.Open(viper.GetString("env.dbFileName"))
	if err != nil {
		fmt.Println("FATAL - Unable to open/find/access database. Exiting application...")
		log.Fatal(err.Error())
	}

	log.Println("DB Opened at:", ArkTestDB.Path)
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
	viper.SetDefault("env.txIterations", 1)
	viper.SetDefault("env.dbfilename", "db/testlog.db")
}

func checkTx4Confirmations() {

}

func main() {
	arkapi := core.NewArkClient(nil)
	arkapi = arkapi.SetActiveConfiguration(core.DEVNET)

	t0 := time.Now()

	for xx := 0; xx < viper.GetInt("env.txIterations"); xx++ {
		arkapi = arkapi.SetActiveConfiguration(core.DEVNET)

		var payload core.TransactionPayload

		for i := 0; i < viper.GetInt("env.txPerPayload"); i++ {
			tx := core.CreateTransaction(viper.GetString("account.recepient"),
				1,
				viper.GetString("env.txDescription"),
				viper.GetString("account.passphrase"), viper.GetString("account.secondPassphrase"))
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
