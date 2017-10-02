package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
