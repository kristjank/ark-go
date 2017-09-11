package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kristjank/ark-go/cmd/arkgoserver/api"
	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var router *gin.Engine

func init() {
	initLogger()
	loadConfig()
	api.InitGlobals()
}

func initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("log/arkgo-server.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(file))
	} else {
		log.Error("Failed to log to file, using default stderr")
	}

}

func loadConfig() {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("cfg")      // path to look for the config file in
	viper.AddConfigPath("settings") // path to look for the config file in

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

//CORSMiddleware function enabling CORS requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func initializeRoutes() {
	log.Info("Initializing routes")

	router.Use(CORSMiddleware())
	// Group peer related routes together
	peerRoutes := router.Group("/voters")
	peerRoutes.Use(api.CheckServiceModelHandler())
	{
		peerRoutes.GET("/rewards", api.GetVoters)
		peerRoutes.GET("/blocked", api.GetBlocked)
	}
	deleRoutes := router.Group("/delegate")
	deleRoutes.Use(api.CheckServiceModelHandler())
	{
		deleRoutes.GET("", api.GetDelegate)
		deleRoutes.GET("/config", api.GetDelegateSharingConfig)
		deleRoutes.GET("/paymentruns", api.GetDelegatePaymentRecord)
		deleRoutes.GET("/paymentruns/details", api.GetDelegatePaymentRecordDetails)
	}
	serviceRoutes := router.Group("/service")
	serviceRoutes.Use(api.OnlyLocalCallAllowed())
	{
		serviceRoutes.GET("/start", api.EnterServiceMode)
		serviceRoutes.GET("/stop", api.LeaveServiceMode)
	}
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("cfg/banner.txt")
	fmt.Print(string(dat))
}

///////////////////////////
func main() {
	printBanner()
	log.Info("..........GOARK-DELEGATE-POOL-SERVER-STARTING............")

	//sending ARKGO Server that we are working with payments
	//setting the version
	api.ArkGoServerVersion = "v0.2.0"

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	pNodeInfo := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
	router.Run(pNodeInfo)

}
