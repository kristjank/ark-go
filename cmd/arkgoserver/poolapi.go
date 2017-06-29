package main

import (
	"strings"

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

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"))

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

//GetBlocked Returns a list of peers to client call. Response is in JSON
func GetBlocked(c *gin.Context) {
	blockedList := viper.GetString("voters.blocklist")

	c.JSON(200, gin.H{
		"blockedList": strings.Split(blockedList, ",")})
}

//GetDelegateSharingConfig Returns a list of peers to client call. Response is in JSON
func GetDelegateSharingConfig(c *gin.Context) {
	blockedList := viper.GetString("voters.blocklist")

	c.JSON(200, gin.H{
		"shareratio":    viper.GetFloat64("voters.shareratio"),
		"fidelity":      viper.GetBool("voters.fidelity"),
		"fidelityLimit": viper.GetInt("voters.fidelityLimit"),
		"minamount":     viper.GetInt("voters.minamount"),
		"deductTxFees":  viper.GetBool("voters.deductTxFeed"),
		"blockedList":   strings.Split(blockedList, ",")})
}
