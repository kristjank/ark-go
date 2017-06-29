package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

var errorlog *os.File
var logger *log.Logger
var router *gin.Engine

func init() {
	initLogger()
	loadConfig()
}

func initLogger() {
	errorlog, err := os.OpenFile("log/goark-node.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.New(errorlog, "ark-go: ", log.Lshortfile|log.LstdFlags)
}

func loadConfig() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("cfg")    // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file

	if err != nil {
		logger.Println("No productive config found - loading sample")
		// try to load sample config
		viper.SetConfigName("sample.config")
		viper.AddConfigPath("cfg")
		err := viper.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			logger.Println("No configuration file loaded - using defaults")
		}
	}

	viper.SetDefault("delegate.address", "")
	viper.SetDefault("delegate.pubkey", "")
	viper.SetDefault("delegate.Daddress", "")
	viper.SetDefault("delegate.Dpubkey", "")

	viper.SetDefault("voters.shareRatio", 0.0)
	viper.SetDefault("voters.txdescription", "share tx by ark-go")
	viper.SetDefault("voters.fidelity", true)
	viper.SetDefault("voters.fidelityLimit", 24)
	viper.SetDefault("voters.minamount", 0.0)
	viper.SetDefault("voters.deductTxFees", true)

	viper.SetDefault("costs.address", "")
	viper.SetDefault("costs.shareRatio", 0.0)
	viper.SetDefault("costs.txdescription", "cost tx by ark-go")
	viper.SetDefault("costs.Daddress", "")

	viper.SetDefault("reserve.address", "")
	viper.SetDefault("reserve.shareRatio", 0.0)
	viper.SetDefault("reserve.txdescription", "reserve tx by ark-go")
	viper.SetDefault("reserve.Daddress", "")

	viper.SetDefault("personal.address", "")
	viper.SetDefault("personal.shareRatio", 0.0)
	viper.SetDefault("personal.txdescription", "personal tx by ark-go")
	viper.SetDefault("personal.Daddress", "")

	viper.SetDefault("client.network", "DEVNET")

	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 54000)
	viper.SetDefault("server.version", "0.1.0")
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("cfg/banner.txt")
	fmt.Print(string(dat))
}

///////////////////////////
func main() {
	printBanner()
	logger.Println("GOARK-DELEGATE-POOL-SERVER-STARTING")

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	pNodeInfo := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
	router.Run(pNodeInfo)

}
