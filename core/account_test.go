package core

import (
	"log"
	"testing"
)

func TestGetAccount(t *testing.T) {
	arkapi := NewArkClient(nil)
	params := AccountQueryParams{Address: "ASJBHz4JfWVUGDyN61hMMnW1Y4ZCTBHL1K"}

	accRest, resp, err := arkapi.GetAccount(params)
	if accRest.Success {
		log.Println(t.Name(), "Account found", accRest.Account.PublicKey, accRest.Account.Balance)

	} else {
		t.Error("Account not found", resp.Status, err.Error())
	}
}
