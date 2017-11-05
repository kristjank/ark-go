package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func createTestRecord() *TestLogRecord {
	return &TestLogRecord{
		TestStarted:        time.Now(),
		TxPerPayload:       viper.GetInt("env.txPerPayload"),
		TxIterations:       viper.GetInt("env.txIterations"),
		TxMultiBroadCast:   viper.GetInt("env.txMultiBroadCast"),
		TxDescription:      viper.GetString("env.txDescription"),
		CreatedAt:          time.Now(),
		ArkGoTesterVersion: ArkGoTesterVersion,
	}
}

func createTestIterationRecord(parentID int) *TestLogIteration {
	return &TestLogIteration{
		TestLogRecordID:  parentID,
		IterationStarted: time.Now(),
	}
}

//Update or Save the record
func (tr *TestLogRecord) Update() {
	err := ArkTestDB.Update(tr)
	if err != nil {
		log.Error("Error updating TestLogRecord: ", err.Error())
	}
}

//Save the record
func (tr *TestLogRecord) Save() {
	err := ArkTestDB.Save(tr)
	if err != nil {
		log.Error("Error saving TestLogRecord: ", err.Error())
	}
}

//Update or Save the record
func (tr *TestLogIteration) Update() {
	err := ArkTestDB.Update(tr)
	if err != nil {
		log.Error("Error updating TestLogIteration: ", err.Error())
	}
}

//Save the record
func (tr *TestLogIteration) Save() {
	err := ArkTestDB.Save(tr)
	if err != nil {
		log.Error("Error saving TestLogIteration: ", err.Error())
	}
}

func listTestRecordsDB() {
	var results []TestLogRecord
	err := ArkTestDB.All(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, element := range results {
		fmt.Println(element)
	}
}

func listTestIterationsRecordsDB() {
	var results []TestLogIteration
	err := ArkTestDB.All(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, element := range results {
		fmt.Println(element)
	}
}
