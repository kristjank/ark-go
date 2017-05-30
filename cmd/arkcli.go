package cmd

import (
	"ark-go/core"
	"fmt"
	"log"
)

func calculcateVoteRatio() {
	arkapi := core.NewArkClient(nil)
	params := core.DelegateQueryParams{PublicKey: "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"}

	votersEarnings := arkapi.CalculateVotersProfit(params, 0.70)

	sumEarned := 0
	sumRatio := 0.0
	sumShareEarned := 0
	for _, element := range votersEarnings {
		log.Println(fmt.Sprintf("|%s|%20d|%15.10f|%15d|%15d|%4d|",
			element.Address,
			element.VoteWeight,
			element.VoteWeightShare,
			element.EarnedAmount100,
			element.EarnedAmountXX,
			element.VoteDuration))

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}
	log.Println("Full forged amount: ", sumEarned, "Ratio calc check sum: ", sumRatio, "Amount to voters: ", sumShareEarned, "Ratio shared: ", float64(sumShareEarned)/float64(sumEarned))

}
