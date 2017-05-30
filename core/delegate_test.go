package core

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
)

func TestListDelegates(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := DelegateQueryParams{Offset: 100}

	deleResponse, _, err := arkapi.ListDelegates(params)
	if deleResponse.Success {
		log.Println(t.Name(), "Success, returned ", deleResponse.TotalCount, "delegates, received:", len(deleResponse.Delegates))
		/*for _, element := range deleResponse.Delegates {
			log.Println(element.Username)
		}*/
	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegateUsername(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := DelegateQueryParams{UserName: "acf"}

	deleResponse, _, err := arkapi.GetDelegate(params)
	if deleResponse.Success {

		out, _ := json.Marshal(deleResponse.SingleDelegate)
		log.Println(t.Name(), "Success, returned", string(out))

	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegatePubKey(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := DelegateQueryParams{PublicKey: "03e6397071866c994c519f114a9e7957d8e6f06abc2ca34dc9a96b82f7166c2bf9"}

	deleResponse, _, err := arkapi.GetDelegate(params)
	if deleResponse.Success {

		out, _ := json.Marshal(deleResponse.SingleDelegate)
		log.Println(t.Name(), "Success, returned", string(out))

	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegateVoters(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := DelegateQueryParams{PublicKey: "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"}

	deleResponse, _, err := arkapi.GetDelegateVoters(params)
	if deleResponse.Success {

		//calculating vote weight
		balance := 0
		for _, element := range deleResponse.Accounts {
			intBalance, _ := strconv.Atoi(element.Balance)
			balance += intBalance
		}

		log.Println(t.Name(), "Success, returned", len(deleResponse.Accounts), "voters for delegate with weight", balance)

	} else {
		t.Error(err.Error())
	}
}

func TestGetDelegateVoteWeight(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := DelegateQueryParams{PublicKey: "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"}

	voteWeight, _, _ := arkapi.GetDelegateVoteWeight(params)

	log.Println(t.Name(), "Success, returned delegate vote weight is", voteWeight)
}

func TestCalculcateVotersProfit(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := DelegateQueryParams{PublicKey: "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"}

	delegateRes, _, _ := arkapi.GetDelegate(params)
	voters, _, _ := arkapi.GetDelegateVoters(params)
	accountRes, _, _ := arkapi.GetAccount(AccountQueryParams{Address: delegateRes.SingleDelegate.Address})

	votersEarnings := CalculateVotersProfit(voters, delegateRes.SingleDelegate, accountRes.Account)

	//log.Println(t.Name(), "Success", votersEarnings)

	for _, element := range votersEarnings {

		log.Println(fmt.Sprintf("|%s|%30d|%30d|%30d|", element.Address, element.VoteWeight, element.VoteWeightShare, element.EarnedAmmount))
	}
}
