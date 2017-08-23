package main

import (
	"fmt"
	"time"

	"github.com/asdine/storm"
	"github.com/kristjank/ark-go/cmd/model"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func save2db(ve core.DelegateDataProfit, tx *core.Transaction, relID int) {
	dbData := model.PaymentLogRecord{}

	dbData.Address = ve.Address
	dbData.VoteWeight = ve.VoteWeight
	dbData.VoteWeightShare = ve.VoteWeightShare
	dbData.EarnedAmount100 = ve.EarnedAmount100
	dbData.EarnedAmountXX = ve.EarnedAmountXX
	dbData.VoteDuration = ve.VoteDuration
	dbData.Transaction = *tx
	dbData.PaymentRecordID = relID
	dbData.CreatedAt = time.Now()

	err := arkpooldb.Save(&dbData)
	if err != nil {
		log.Error(err.Error())
	}
}

func listPaymentsDetailsFromDB() {
	var results []model.PaymentLogRecord
	err := arkpooldb.All(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, element := range results {
		fmt.Println(element)
	}
}

func listPaymentsDB() {
	var results []model.PaymentRecord
	err := arkpooldb.All(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, element := range results {
		fmt.Println(element)
	}
}

func initializeBoltClient() {
	var err error
	arkpooldb, err = storm.Open(viper.GetString("client.dbfilename"))

	if err != nil {
		log.Panic(err.Error())
	}

	log.Println("DB Opened at:", arkpooldb.Path)
	//defer arkpooldb.Close()
}

func createPaymentRecord() model.PaymentRecord {

	delegateAddress := viper.GetString("delegate.address")
	if viper.GetString("client.network") == "DEVNET" {
		delegateAddress = viper.GetString("delegate.Daddress")
	}

	payRec := model.PaymentRecord{
		ShareRatio:    viper.GetFloat64("voters.shareratio"),
		CostsRatio:    viper.GetFloat64("costs.shareratio"),
		PersonalRatio: viper.GetFloat64("personal.shareratio"),
		ReserveRatio:  viper.GetFloat64("reserve.shareratio"),
		CreatedAt:     time.Now(),
		FeeDeduction:  viper.GetBool("voters.deductTxFees"),
		Fidelity:      viper.GetBool("voters.fidelity"),
		FidelityLimit: viper.GetInt("voters.fidelityLimit"),
		MinAmount:     viper.GetFloat64("voters.minamount"),
		Delegate:      delegateAddress,
	}
	return payRec
}
