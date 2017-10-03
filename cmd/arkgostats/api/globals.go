package api

import (
	"fmt"

	"github.com/asdine/storm"

	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ArkAPIclient *core.ArkClient
var ArkStatsDB *storm.DB
var ArkGoStatsServerVersion string

func InitGlobals() {
	ArkAPIclient = core.NewArkClient(nil)
	openDB()
}

func openDB() {
	log.Info("Opening/Reopening database")
	var err error
	ArkStatsDB, err = storm.Open(viper.GetString("server.dbfilename"))
	if err != nil {
		fmt.Println("FATAL - Unable to open/find/access database. Exiting application...")
		log.Fatal(err.Error())
	}

	log.Println("DB Opened at:", ArkStatsDB.Path)
}

func closeDB() {
	log.Info("Closing database")
	err := ArkStatsDB.Close()
	if err != nil {
		log.Error(err.Error())
	}
}
