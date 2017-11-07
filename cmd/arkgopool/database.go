package main

import (
	"fmt"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/kristjank/ark-go/cmd/model"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func beginTx() storm.Node {
	dbtx, err := arkpooldb.Begin(true)
	if err != nil {
		log.Error(err.Error())
	}
	return dbtx
}

func commitTx(dbtx storm.Node) {
	err := dbtx.Commit()
	if err != nil {
		log.Error(err.Error())
	}
}

func rollbackTx(dbtx storm.Node) {
	err := dbtx.Rollback()
	if err != nil {
		log.Error(err.Error())
	}
}

func save2db(dbtx storm.Node, ve core.DelegateDataProfit, tx *core.Transaction, relID int) {
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

	err := dbtx.Save(&dbData)
	if err != nil {
		log.Error(err.Error())
	}
}

func savebonus2db(dbtx storm.Node, address string, tx *core.Transaction, relID int) {
	dbData := model.PaymentLogRecord{}

	dbData.Address = address
	dbData.Transaction = *tx
	dbData.PaymentRecordID = relID
	dbData.CreatedAt = time.Now()

	err := dbtx.Save(&dbData)
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

func getTxIDsFromPaymentLogRecord(payRec model.PaymentRecord) ([]string, error) {
	var results []model.PaymentLogRecord
	var err error
	var query storm.Query
	var transIDList []string

	query = arkpooldb.Select(q.Eq("PaymentRecordID", payRec.Pk))
	err = query.Find(&results)

	if err != nil {
		log.Error(err.Error())
		return transIDList, err
	}

	for _, el := range results {
		transIDList = append(transIDList, el.Transaction.ID)
	}
	return transIDList, nil
}

func getLastPaymentRecord() (model.PaymentRecord, error) {
	var results []model.PaymentRecord
	err := arkpooldb.All(&results, storm.Limit(1), storm.Reverse())

	rec := model.PaymentRecord{}

	if err != nil {
		log.Error(err.Error())
		return rec, err
	}
	rec = results[0]
	return rec, nil
}

func createPaymentRecord() model.PaymentRecord {
	delegateAddress := viper.GetString("delegate.address")
	delegatePubKey := viper.GetString("delegate.pubkey")
	if viper.GetString("client.network") == "DEVNET" {
		delegateAddress = viper.GetString("delegate.Daddress")
		delegatePubKey = viper.GetString("delegate.Dpubkey")
	}

	payRec := model.PaymentRecord{
		ShareRatio:       viper.GetFloat64("voters.shareratio"),
		CostsRatio:       viper.GetFloat64("costs.shareratio"),
		PersonalRatio:    viper.GetFloat64("personal.shareratio"),
		ReserveRatio:     viper.GetFloat64("reserve.shareratio"),
		CreatedAt:        time.Now().UTC(),
		FeeDeduction:     viper.GetBool("voters.deductTxFees"),
		Fidelity:         viper.GetBool("voters.fidelity"),
		FidelityLimit:    viper.GetInt("voters.fidelityLimit"),
		MinAmount:        viper.GetFloat64("voters.minamount"),
		Network:          viper.GetString("client.network"),
		Blocklist:        viper.GetString("voters.blocklist"),
		Whitelist:        viper.GetString("voters.whitelist"),
		CapBalance:       viper.GetBool("voters.capBalance"),
		BalanceCapAmount: viper.GetFloat64("voters.balanceCapAmount"),
		BlockBalanceCap:  viper.GetBool("voters.blockBalanceCap"),
		Delegate:         delegateAddress,
		DelegatePubKey:   delegatePubKey,
		ArkGoPoolVersion: ArkGoPoolVersion,
	}
	return payRec
}
