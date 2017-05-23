package core

import (
	"log"
	"testing"
)

func TestListDelegates(t *testing.T) {
	arkapi := NewArkClient(nil)

	//params := TransactionQueryParams{Limit: 10, SenderID: "AQLUKKKyKq5wZX7rCh4HJ4YFQ8bpTpPJgK"}

	deleResponse, _, err := arkapi.ListDelegates()
	if deleResponse.Success {
		log.Println(t.Name(), "Success, returned ", deleResponse.TotalCount, "delegates")
	} else {
		t.Error(err.Error())
	}
}
