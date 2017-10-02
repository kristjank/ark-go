package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kristjank/ark-go/cmd/model"
	log "github.com/sirupsen/logrus"
)

func OnlyLocalCallAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.ClientIP() == "127.0.0.1" || c.ClientIP() == "::1" {
			c.Next()
		} else {
			log.Info("Outside call to service mode is not allowed")
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
}

//GetServerInformation Returns a server statistics
func GetServerInformation(c *gin.Context) {
	c.JSON(200, gin.H{"version": ArkGoStatsServerVersion})
}

//ReceivePaymetLog from blockchain
func ReceivePaymetLog(c *gin.Context) {
	var recv model.PaymentRecord
	err := c.BindJSON(&recv)
	if err != nil {
		log.Error(err.Error())
	}

	recv.SourceIP = c.ClientIP()

	err = ArkStatsDB.Save(&recv)
	c.JSON(200, gin.H{"success": true, "logID": recv.Pk})
}
