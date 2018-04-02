package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/fatih/color"
	"github.com/kristjank/ark-go/cmd/arkgostats/api"
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
	file, err := os.OpenFile("log/arkgo-stats.log", os.O_CREATE|os.O_WRONLY, 0666)
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

	viper.SetDefault("server.address", "0.0.0.0")
	viper.SetDefault("server.port", 54010)
	viper.SetDefault("server.dbfilename", "db/arkstats.db")
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
	statsRoutes := router.Group("/")
	statsRoutes.Use()
	{
		//statsRoutes.PUT("/payment")
		statsRoutes.GET("info", api.GetServerInformation)
		statsRoutes.POST("log/payment", api.ReceivePaymetLog)
		statsRoutes.GET("/payments", api.SendPaymentLog)
		statsRoutes.GET("/delegate/:address", api.SendPaymentLog4Delegate)
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
	log.Info("..........ARKGO-STATS-SERVER-STARTING............")

	//sending ARKGO Server that we are working with payments
	//setting the version
	api.ArkGoStatsServerVersion = "v0.1.2"

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	pNodeInfo := fmt.Sprintf("%s:%d", viper.GetString("server.address"), viper.GetInt("server.port"))
	router.Run(pNodeInfo)
}
