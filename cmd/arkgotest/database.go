package main

import (
	"time"

	"github.com/spf13/viper"
)

func (t *TestLogRecord) createTestRecord() TestLogRecord {
	testRecord := TestLogRecord{
		txPerPayload:       viper.GetInt("env.txPerPayload"),
		txIterations:       viper.GetInt("env.txIterations"),
		txMultiBroadCast:   viper.GetInt("env.txPerPayload"),
		txDescription:      viper.GetString("env.txDescription"),
		CreatedAt:          time.Now(),
		ArkGoTesterVersion: ArkGoTesterVersion,
	}
	return testRecord
}
