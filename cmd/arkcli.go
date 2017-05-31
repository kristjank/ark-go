package cmd

import (
	"ark-go/core"
	"fmt"
	"log"
)

func calculcateVoteRatio() {
	arkapi := core.NewArkClient(nil)

	deleKey := "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		deleKey = "02bcfa0951a92e7876db1fb71996a853b57f996972ed059a950d910f7d541706c9 "
	}

	params := core.DelegateQueryParams{PublicKey: deleKey}

	votersEarnings := arkapi.CalculateVotersProfit(params, 0.70)

	var payload core.TransactionPayload

	//log.Println(t.Name(), "Success", votersEarnings)
	fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)

	sumEarned := 0.9
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := float64(len(votersEarnings)) * (float64(core.EnvironmentParams.Fees.Send) / core.SATOSHI)
	for _, element := range votersEarnings {
		log.Println(fmt.Sprintf("|%s|%15.8f|%15.8f|%15.8f|%15.8f|%4d|%25d|",
			element.Address,
			element.VoteWeight,
			element.VoteWeightShare,
			element.EarnedAmount100,
			element.EarnedAmountXX,
			element.VoteDuration,
			int(element.EarnedAmountXX*core.SATOSHI)))

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		//transaction parameters
		tx := core.CreateTransaction(element.Address,
			int64(element.EarnedAmountXX*core.SATOSHI),
			"chris: 1st profit sharing payment... |tx made with ark-go",
			"",
			"")

		payload.Transactions = append(payload.Transactions, tx)

	}
	log.Println("Full forged amount: ", sumEarned, "Ratio calc check sum: ", sumRatio, "Amount to voters: ", sumShareEarned, "Ratio shared: ", float64(sumShareEarned)/float64(sumEarned), "Lottery:", int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI))
	log.Println(fmt.Sprintf("Payment fees: %2.2f", feeAmount))

	tx := core.CreateTransaction("ANqeL7CP2som7q9NFbRuaUc5WUnwYkSbFY",
		int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI),
		"chris: 1st month lottery fund reserve... |tx made with ark-go",
		"",
		"")

	payload.Transactions = append(payload.Transactions, tx)

	/*//payload complete - posting
	res, httpresponse, err := arkapi.PostTransaction(payload)
	if res.Success {
		log.Println("Success,", httpresponse.Status, res.TransactionIDs)

	} else {
		log.Println(res.Message, res.Error, httpresponse.Status, err.Error())

	}*/
}
