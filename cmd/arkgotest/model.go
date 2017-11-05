package main

import (
	"time"
)

//TestLogRecord structure
type TestLogRecord struct {
	ID                 int `storm:"id,increment"` // primary key with auto increment
	TxPerPayload       int
	TxIterations       int
	TxMultiBroadCast   int
	TxDescription      string
	TestStarted        time.Time
	TestStopped        time.Time
	TestStatus         string
	TestLogIterationID int       `storm:"index"`
	CreatedAt          time.Time `storm:"index"`
	ArkGoTesterVersion string
}

//TestLogIteration structure
type TestLogIteration struct {
	ID               int `storm:"id,increment"`
	IterationStarted time.Time
	IterationStopped time.Time
	TestStatus       string
	TxIDs            []string `storm:"index"`
	TestLogRecordID  int      `storm:"index"`
}

//TestLogTransaction structure
type TestLogTransaction struct {
	ID              int    `storm:"id,increment"`
	TransactionID   string `storm:"index"`
	Confirmations   int
	TestLogRecordID int       `storm:"index"`
	CreatedAt       time.Time `storm:"index"`
}
