package model

import (
	"time"

	"github.com/kristjank/ark-go/core"
)

//PaymentLogRecord structure
//Child structure
type PaymentLogRecord struct {
	Pk              int    `storm:"id,increment"` // primary key with auto increment
	Address         string `storm:"index"`
	VoteWeight      float64
	VoteWeightShare float64
	EarnedAmount100 float64
	EarnedAmountXX  float64
	VoteDuration    int
	Transaction     core.Transaction
	PaymentRecordID int       `storm:"index"`
	CreatedAt       time.Time `storm:"index"`
}

//PaymentRecord structure
//MainStructure
type PaymentRecord struct {
	Pk               int    `storm:"id,increment"`
	Delegate         string `storm:"index"`
	ShareRatio       float64
	CostsRatio       float64
	ReserveRatio     float64
	PersonalRatio    float64
	Fidelity         bool
	FidelityLimit    int
	MinAmount        float64
	FeeDeduction     bool
	FeeAmount        float64
	NrOfTransactions int
	VoteWeight       int
	CreatedAt        time.Time `storm:"index"`
}
