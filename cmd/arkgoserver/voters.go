package main

import (
	"github.com/kristjank/ark-go/core"
	"github.com/spf13/viper"
	"gopkg.in/gin-gonic/gin.v1"
)

var arkclient = core.NewArkClient(nil)

//GetVoters Returns a list of peers to client call. Response is in JSON
func GetVoters(c *gin.Context) {

	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"))

	c.JSON(200, votersEarnings)
}

//GetDelegate Returns a list of peers to client call. Response is in JSON
func GetDelegate(c *gin.Context) {

	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}

	deleResp, _, _ := arkclient.GetDelegate(params)

	c.JSON(200, deleResp)
}
