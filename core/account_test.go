package core

import (
	"log"
	"testing"
)

func TestGetAccount(t *testing.T) {
	arkapi := NewArkClient(nil)
	params := AccountQueryParams{Address: "ANqeL7CP2som7q9NFbRuaUc5WUnwYkSbFY"}

	accRest, resp, err := arkapi.GetAccount(params)
	if accRest.Success {
		log.Println(t.Name(), "Account found", accRest.Account.PublicKey, accRest.Account.Balance, accRest.Account.SecondSignature)

	} else {
		t.Error("Account not found", resp.Status, err.Error())
	}
}
