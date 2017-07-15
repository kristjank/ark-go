package main

import (
	"gopkg.in/gin-gonic/gin.v1"
)

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
	logger.Println("Initializing routes")

	router.Use(CORSMiddleware())
	// Group peer related routes together
	peerRoutes := router.Group("/voters")
	{
		peerRoutes.GET("/rewards", GetVoters)
		peerRoutes.GET("/blocked", GetBlocked)
	}
	deleRoutes := router.Group("/delegate")
	{
		deleRoutes.GET("", GetDelegate)
		deleRoutes.GET("/config", GetDelegateSharingConfig)
		deleRoutes.GET("/paymentruns", GetDelegatePaymentRecord)
		deleRoutes.GET("/paymentruns/details", GetDelegatePaymentRecordDetails)
	}
}
