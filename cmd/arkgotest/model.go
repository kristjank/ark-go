package main

import (
	"time"
)

//TestLogRecord structure
type TestLogRecord struct {
	Pk                 int `storm:"id,increment,index"` // primary key with auto increment
	txPerPayload       string
	txIterations       string
	txMultiBroadCast   string
	txDescription      string
	TestStarted        time.Time
	TestStopped        time.Time
	TestStatus         string
	TestLogIterationID int       `storm:"index"`
	CreatedAt          time.Time `storm:"index"`
}

//TestLogIteration structure
type TestLogIteration struct {
	Pk                     int `storm:"id,increment,index"`
	IterationStarted       time.Time
	IterationStopped       time.Time
	TestStatus             string
	IterationTransactionID int       `storm:"index"`
	CreatedAt              time.Time `storm:"index"`
}

//IterationTransaction structure
type IterationTransaction struct {
	Pk            int    `storm:"id,increment,index"`
	TransactionID string `storm:"index"`
	Confirmations int
	CreatedAt     time.Time `storm:"index"`
}
