package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/asdine/storm"
	"github.com/fatih/color"
	"github.com/kristjank/ark-go/core"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var ArkGoTesterVersion string

//ArkAPIClient - ARKAPI Client
var ArkAPIClient *core.ArkClient

//ArkTestDB - testDB to store log records
var ArkTestDB *storm.DB

//ConsoleReader console input reader
var ConsoleReader = bufio.NewReader(os.Stdin)

var wg sync.WaitGroup

func init() {
	initLogger()
	loadConfig()
	ArkAPIClient = core.NewArkClient(nil)
	ArkAPIClient = ArkAPIClient.SetActiveConfiguration(core.DEVNET)
	openDB()
	dumpConfig()
}

func openDB() {
	log.Info("Opening/Reopening database")
	var err error
	ArkTestDB, err = storm.Open(viper.GetString("env.dbFileName"))
	if err != nil {
		fmt.Println("FATAL - Unable to open/find/access database. Exiting application...")
		log.Fatal(err.Error())
	}

	log.Println("DB Opened at:", ArkTestDB.Path)
}

func initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("log/arkgotest.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(file))
	} else {
		log.Error("Failed to log to file, using default stderr")
	}
}

func loadConfig() {
	viper.AddConfigPath("cfg")  // path to look for the config file in
	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		log.Info("No productive config found - loading sample")
		// try to load sample config
		viper.SetConfigName("sample.config")
		viper.AddConfigPath("cfg")
		err := viper.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			log.Fatal("No configuration file loaded - using defaults")
		}
	}

	viper.SetDefault("env.txPerPayload", 250)
	viper.SetDefault("env.txMultiBroadCast", 1)
	viper.SetDefault("env.txIterations", 1)
	viper.SetDefault("env.dbfilename", "db/testlog.db")
	viper.SetDefault("env.txDescription", "ARK Test tx script")
	viper.SetDefault("env.singlePeerTest", false)
	viper.SetDefault("env.singlePeerIp", "")
}

func dumpConfig() {
	log.Info("--- Running ArkGO Test with the following config ---")
	log.Info("Iterations: ", viper.GetInt("env.txIterations"))
	log.Info("Tx per payload: ", viper.GetInt("env.txPerPayload"))
	log.Info("Recepient of tx: ", viper.GetString("account.recepient"))
	log.Info("Tx description: ", viper.GetString("env.txDescription"))
	log.Info("--- end of config params lisitngs ---")
}

func main() {
	ArkGoTesterVersion = "v0.1.0"

	log.Info("=============================================================================")
	log.Info("ARKGO Tester application starting")
	log.Info("ArkApiClient connected, active peer: ", ArkAPIClient.GetActivePeer())

	//test code
	/*	lastRec, err := getLatstTestRecord()
		if err == nil {
			findConfirmations(lastRec)
			checkConfirmations(lastRec)
		}
	*/
	//SILENT MODE CHECKING AND AUTOMATION RUNNING
	modeSilentPtr := flag.Bool("silent", false, "Is silent mode")
	//autoPayment := flag.Bool("autopay", true, "Process auto payment")
	flag.Parse()
	log.Info(flag.Args())
	if *modeSilentPtr {
		log.Info("Silent Mode active")
		log.Info("Starting to send Test trx")
		runTests()
		log.Info("Exiting silent mode")
		os.Exit(1985)
	}

	var choice = 1
	for choice != 0 {
		//pause()
		printMenu()

		//fmt.Scan(&choice)
		fmt.Fscan(ConsoleReader, &choice)
		ConsoleReader.ReadString('\n')

		switch choice {
		case 1:
			clearScreen()
			color.Set(color.FgMagenta)
			runTests()
			color.Unset()
		case 8:
			clearScreen()
			color.Set(color.FgHiWhite)
			lastRec, err := getLatstTestRecord()
			if err == nil {
				findConfirmations(lastRec)
				checkConfirmations(lastRec)
			}
			pause()
			color.Unset()
		case 9:
			clearScreen()
			color.Set(color.FgHiGreen)
			listTestRecordsDB()
			//listTestIterationsRecordsDB()
			pause()
			color.Unset()
		}
	}
	color.Unset()
	log.Info("Exiting ARKGOTest Application....")
}
