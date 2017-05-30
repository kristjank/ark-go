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

	votersEarnings := arkapi.CalculateVotersProfit(params, 0.70)

	//log.Println(t.Name(), "Success", votersEarnings)
	//log.Println(t.Name(), "Success", votersEarnings)
	sumEarned := 0.9
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := float64(len(votersEarnings)) * (float64(EnvironmentParams.Fees.Send) / SATOSHI)
	for _, element := range votersEarnings {
		log.Println(fmt.Sprintf("|%s|%15.8f|%15.8f|%15.8f|%15.8f|%4d|%25d|",
			element.Address,
			element.VoteWeight,
			element.VoteWeightShare,
			element.EarnedAmount100,
			element.EarnedAmountXX,
			element.VoteDuration,
			int(element.EarnedAmountXX*SATOSHI)))

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}
	log.Println("Full forged amount: ", sumEarned, "Ratio calc check sum: ", sumRatio, "Amount to voters: ", sumShareEarned, "Ratio shared: ", float64(sumShareEarned)/float64(sumEarned), "Lottery:", int64((sumEarned-sumShareEarned-feeAmount)*SATOSHI))
	log.Println(fmt.Sprintf("Payment fees: %2.2f", feeAmount))
}
