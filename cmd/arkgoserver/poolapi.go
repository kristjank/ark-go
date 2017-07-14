package main

import (
	"strconv"
	"strings"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/kristjank/ark-go/cmd/model"
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

//GetDelegatePaymentRecord Returns a list of peers to client call. Response is in JSON
func GetDelegatePaymentRecord(c *gin.Context) {
	var results []model.PaymentRecord

	err := arkpooldb.All(&results)

	if err == nil {
		c.JSON(200, results)
	} else {
		c.JSON(200, gin.H{"success": false, "error": err.Error()})
	}
	c.JSON(200, gin.H{"success": true})

}

//GetDelegatePaymentRecordDetails Returns a list of peers to client call. Response is in JSON
//URL samples:
//1.TO GET ALL PAYMENT DETAILS: http://localhost:54000/delegate/paymentruns/details
//2.TO GET ALL PAYMENT DETAILS FOR SPECIFIED PAYMENT RUN: http://localhost:54000/delegate/paymentruns/details?parentid=12
//3.TO GET ALL PAYMENT DETAILS FOR SPECIFIED VOTER(address): http://ocalhost:54000/delegate/payments/details?address=D5St8ot3asrxYW3o63EV3bM1VC6UBKMUfE
//4.TO GET ALL PAYMENT DETAILS FOR SPECIFIED VOTER(address) in Specified RUN:http://localhost:54000/delegate/paymentruns/details?parentid=12&address=D5St8ot3asrxYW3o63EV3bM1VC6UBKMUfE
func GetDelegatePaymentRecordDetails(c *gin.Context) {
	var results []model.PaymentLogRecord
	var err error
	var query storm.Query

	id, err := strconv.Atoi(c.DefaultQuery("parentid", "-1"))
	address := c.DefaultQuery("address", "")

	if id != -1 && address != "" {
		query = arkpooldb.Select(q.Eq("PaymentRecordID", id), q.Eq("Address", address)).Reverse()
	} else if id != -1 && address == "" {
		query = arkpooldb.Select(q.Eq("PaymentRecordID", id)).Reverse()
	} else if id == -1 && address != "" {
		query = arkpooldb.Select(q.Eq("Address", address)).Reverse()
	} else {
		query = arkpooldb.Select()
	}

	err = query.Find(&results)

	if err == nil {
		c.JSON(200, results)
	} else {
		c.JSON(200, gin.H{"success": false, "error": err.Error()})
	}
}
