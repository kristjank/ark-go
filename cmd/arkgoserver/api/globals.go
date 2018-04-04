package api

import (
	"fmt"
	"sync"
	"time"

	"github.com/asdine/storm"

	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ArkAPIclient *core.ArkClient
var Arkpooldb *storm.DB
var ArkGoServerVersion string

var syncMutex = &sync.RWMutex{}
var voterMutex = &sync.RWMutex{}
var isServiceMode bool
var rewardTicker *time.Ticker
var VotersEarnings []core.DelegateDataProfit

func InitGlobals() {
	isServiceMode = false
	if len(viper.GetString("server.autoconfigPeer")) > 0 {
		log.Info("ARKGOServer client setting properties via autocofig peer ", viper.GetString("server.autoconfigPeer"))
		ArkAPIclient = ArkAPIclient.SetActiveConfigurationFromPeerAddress(viper.GetString("server.autoconfigPeer"))
	} else if viper.GetString("server.network") == "DEVNET" {
		ArkAPIclient = ArkAPIclient.SetActiveConfiguration(core.DEVNET)
	} else if viper.GetString("server.network") == "KAPU" {
		ArkAPIclient = ArkAPIclient.SetActiveConfiguration(core.KAPU)
	} else {
		ArkAPIclient = ArkAPIclient.SetActiveConfiguration(core.MAINNET)
	}
	openDB()

	initTicker4PendingRewardCalculation()
}

func initTicker4PendingRewardCalculation() {
	rewardTicker = time.NewTicker(time.Minute * 10)
	pubKey := viper.GetString("delegate.pubkey")
	params := core.DelegateQueryParams{PublicKey: pubKey}

	//do first reading of calculations (first tick is in 10 minutes from now)
	go func() {
		voterMutex.Lock()
		VotersEarnings = ArkAPIclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"), viper.GetString("voters.whitelist"), viper.GetBool("voters.capBalance"), viper.GetFloat64("voters.BalanceCapAmount")*core.SATOSHI, viper.GetBool("voters.blockBalanceCap"))
		voterMutex.Unlock()
	}()

	go func() {
		for t := range rewardTicker.C {
			log.Info("Calling voter earning cache calculation for faster display", t)
			fmt.Println("Calling voter earning cache calculation for faster display", t)

			tmpEarnings := ArkAPIclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"), viper.GetString("voters.whitelist"), viper.GetBool("voters.capBalance"), viper.GetFloat64("voters.BalanceCapAmount")*core.SATOSHI, viper.GetBool("voters.blockBalanceCap"))

			voterMutex.Lock()
			VotersEarnings = tmpEarnings
			voterMutex.Unlock()
		}
	}()
}

func openDB() {
	log.Info("Opening/Reopening database")
	var err error
	Arkpooldb, err = storm.Open(viper.GetString("server.dbfilename"))
	if err != nil {
		fmt.Println("FATAL - Unable to open/find/access database. Exiting application...")
		log.Fatal(err.Error())
	}

	log.Println("DB Opened at")
}

func closeDB() {
	log.Info("Closing database")
	err := Arkpooldb.Close()
	if err != nil {
		log.Error(err.Error())
	}
}
