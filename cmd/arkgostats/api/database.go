package api

import (
	"github.com/asdine/storm/q"
	"github.com/kristjank/ark-go/cmd/model"
	log "github.com/sirupsen/logrus"
)

func getPayments(offset int, network string) ([]model.PaymentRecord, error) {
	var results []model.PaymentRecord

	query := ArkStatsDB.Select(q.Eq("Network", network)).Reverse().Limit(50).Skip(offset)
	err := query.Find(&results)

	if err != nil {
		log.Error("getPayments ", err.Error())
	}

	return results, err
}
