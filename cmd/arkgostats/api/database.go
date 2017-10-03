package api

import (
	"github.com/asdine/storm"
	"github.com/kristjank/ark-go/cmd/model"
	log "github.com/sirupsen/logrus"
)

func getPayments(offset int) ([]model.PaymentRecord, error) {
	var results []model.PaymentRecord
	err := ArkStatsDB.AllByIndex("Pk", &results, storm.Limit(50), storm.Skip(offset), storm.Reverse())

	if err != nil {
		log.Error("getPayments ", err.Error())
	}

	return results, err
}
