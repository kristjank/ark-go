package main

import (
	"fmt"

	"github.com/dghubble/sling"
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

func sendStatisticsData(payRec *model.PaymentRecord) {
	response := new(postStatsResponse)
	error := new(postStatsResponse)

	statURL := fmt.Sprintf("http://%s:%d", viper.GetString("client.statPeer"), viper.GetInt("client.statPort"))
	statsBase := sling.New().Base(statURL).Client(nil).Add("Content-Type", "application/json")
	resp, err := statsBase.New().Post("log/payment").BodyJSON(payRec).Receive(response, error)

	if err != nil {
		log.Error("Error sending statistics data", resp, err)
	}
}

func splitAndDeliverPayload(payload core.TransactionPayload) {
	//calculating number of chunks (based on 20tx in one chunk to send to one peer)
	payoutsFolderName := createLogFolder()
	var divided [][]*core.Transaction
	numPeers := len(payload.Transactions) / 20
	if numPeers == 0 {
		numPeers = 1
	}
	chunkSize := (len(payload.Transactions) + numPeers - 1) / numPeers
	if chunkSize == 0 {
		chunkSize = 1
	}

	//sliptting the payload to number of needed peers
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

	}
	if splitcout != len(payload.Transactions) {
		log.Error("TX spliting not OK")
	}
}

func deliverPayloadThreaded(tmpPayload core.TransactionPayload, chunkIx int, logFolder string) {
	numberOfPeers2MultiBroadCastTo := viper.GetInt("client.multibroadcast")
	if numberOfPeers2MultiBroadCastTo > 15 {
		numberOfPeers2MultiBroadCastTo = 15
		log.Warn("Max broadcast number too high - set by user, reseting to value 15")
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
				go log2csv(tmpPayload, res.TransactionIDs, filename, "OK")
			} else {
				color.Set(color.FgHiRed)
				go log2csv(tmpPayload, nil, filename, res.Error)
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
				}

			}
		}(transIDPart, id, arkclient)
	}
	wgConfirmations.Wait()
	log.Info("---------------END OF CONFIRMATION CHECK-----------------")
}
