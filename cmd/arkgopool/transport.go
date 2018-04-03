package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/kristjank/ark-go/cmd/model"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type postStatsResponse struct {
	Success bool   `json:"success,omitempty"`
	LogID   int    `json:"logID,omitempty"`
	Error   string `json:"error,omitempty"`
}

func splitAndDeliverPayload(payload core.TransactionPayload) {
	//calculating number of chunks (based on 20tx in one chunk to send to one peer)
	payoutsFolderName := createLogFolder()
	var divided [][]*core.Transaction

	chunkSize := viper.GetInt("client.payloadsize")
	if chunkSize > 40 {
		chunkSize = 40
	}

	//sliptting the payload to size defined in chunksize
	for i := 0; i < len(payload.Transactions); i += chunkSize {
		end := i + chunkSize
		if end > len(payload.Transactions) {
			end = len(payload.Transactions)
		}
		divided = append(divided, payload.Transactions[i:end])
	}
	//end of spliting transactions

	var tmpPayload core.TransactionPayload
	splitcout := 0
	for chunkIx, h := range divided {
		tmpPayload.Transactions = h
		splitcout += len(h)

		deliverPayloadThreaded(tmpPayload, chunkIx, payoutsFolderName)
		fmt.Println("Sending transactions to the network", strconv.FormatFloat(float64(100*splitcout/len(payload.Transactions)), 'f', 2, 64), "%")

		if splitcout < len(payload.Transactions) {
			time.Sleep(time.Second * 40) //waiting before sending another batch - Quick fix (rewrite after v2 is out)
		}
	}
	if splitcout != len(payload.Transactions) {
		log.Error("TX spliting not OK")
	}
}

func deliverPayloadThreaded(tmpPayload core.TransactionPayload, chunkIx int, logFolder string) {
	numberOfPeers2MultiBroadCastTo := viper.GetInt("client.multibroadcast")
	if numberOfPeers2MultiBroadCastTo > 10 {
		numberOfPeers2MultiBroadCastTo = 10
		log.Warn("Max broadcast number too high - set by user, reseting to value 10")
	}
	log.Info("Starting multibroadcast/multithreaded parallel payout to ", numberOfPeers2MultiBroadCastTo, " number of peers")
	peers := arkclient.GetRandomXPeers(numberOfPeers2MultiBroadCastTo)
	for i := 0; i < numberOfPeers2MultiBroadCastTo; i++ {
		wg.Add(1)

		//treaded function
		go func(tmpPayload core.TransactionPayload, peer core.Peer, chunkIx int, logFolder string) {
			defer wg.Done()
			filename := fmt.Sprintf("log/%s/Batch_%02d_Peer%s.csv", logFolder, chunkIx, peer.IP)

			arkTmpClient := core.NewArkClientFromPeer(peer)
			res, _, _ := arkTmpClient.PostTransaction(tmpPayload)
			if res.Success {
				color.Set(color.FgHiGreen)
				log2csv(tmpPayload, res.TransactionIDs, filename, "OK")
			} else {
				color.Set(color.FgHiRed)
				log2csv(tmpPayload, nil, filename, res.Error)
			}
		}(tmpPayload, peers[i], chunkIx, logFolder)
	}
}

func findConfirmations(payRec model.PaymentRecord) {
	transIDList, err := getTxIDsFromPaymentLogRecord(payRec)
	if err != nil {
		//TODO handle error
		log.Error(err.Error())
		return
	}

	var divided [][]string

	numPeers := len(transIDList) / len(core.EnvironmentParams.Network.PeerList)
	if numPeers == 0 {
		numPeers = 1
	}
	chunkSize := (len(transIDList) + numPeers - 1) / numPeers
	if chunkSize == 0 {
		chunkSize = 1
	}

	//sliptting the payload to number of needed peers
	for i := 0; i < len(transIDList); i += chunkSize {
		end := i + chunkSize
		if end > len(transIDList) {
			end = len(transIDList)
		}
		divided = append(divided, transIDList[i:end])
	}
	//end of spliting transactions

	log.Info("---------------START OF CONFIRMATION CHECK-----------------")
	for id, transIDPart := range divided {
		wgConfirmations.Add(1)

		go func(transIDs []string, idPeer int, arkapi *core.ArkClient) {
			defer wgConfirmations.Done()
			arkTmpClient := core.NewArkClientFromPeer(arkapi.GetRandomXPeers(1)[0])
			for _, txID := range transIDs {
				params := core.TransactionQueryParams{ID: txID}
				arkTransaction, _, _ := arkTmpClient.GetTransaction(params)

				confirmations := 0
				if arkTransaction.Success {
					confirmations = arkTransaction.SingleTransaction.Confirmations
					if confirmations < 1 {
						fmt.Println("Missing transaction ", txID)
					}

				}

			}
		}(transIDPart, id, arkclient)
	}
	wgConfirmations.Wait()
	log.Info("---------------END OF CONFIRMATION CHECK-----------------")
}
