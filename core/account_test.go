package core

import (
	"log"
	"testing"
)

func TestGetAccount(t *testing.T) {
	arkapi := NewArkClient(nil)
	address := "ANqeL7CP2som7q9NFbRuaUc5WUnwYkSbFY"

	if EnvironmentParams.Network.Type == DEVNET {
		address = "DQUjMT6fhJWbwhaYL5pPdX9v5qPiRcAzRb"
	}

	params := AccountQueryParams{Address: address}

	accRest, resp, err := arkapi.GetAccount(params)
	if accRest.Success {
		log.Println(t.Name(), "Account found", accRest.Account.PublicKey, accRest.Account.Balance, accRest.Account.SecondSignature)

	} else {
		if resp != nil {
			t.Error("Account not found", resp.Status, err.Error())
		} else {
			t.Error("Account not found", err.Error())
		}
	}
}
