package api

import "github.com/kristjank/ark-go/cmd/model"

func copyStructure(src, dest *model.PaymentRecord) {
	dest.Delegate = src.Delegate
	dest.DelegatePubKey = src.DelegatePubKey
	dest.ShareRatio = src.ShareRatio
	dest.CostsRatio = src.CostsRatio
	dest.ReserveRatio = src.ReserveRatio
	dest.PersonalRatio = src.PersonalRatio
	dest.Fidelity = src.Fidelity
	dest.FidelityLimit = src.FidelityLimit
	dest.MinAmount = src.MinAmount
	dest.FeeDeduction = src.FeeDeduction
	dest.FeeAmount = src.FeeAmount
	dest.NrOfTransactions = src.NrOfTransactions
	dest.VoteWeight = src.VoteWeight
	dest.Network = src.Network
	dest.Blocklist = src.Blocklist
	dest.Whitelist = src.Whitelist
	dest.CapBalance = src.CapBalance
	dest.BalanceCapAmount = src.BalanceCapAmount
	dest.BlockBalanceCap = src.BlockBalanceCap
	dest.SourceIP = src.SourceIP
	dest.CreatedAt = src.CreatedAt
	dest.ArkGoPoolVersion = src.ArkGoPoolVersion
}
