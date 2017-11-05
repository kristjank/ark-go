package main

import (
	"fmt"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/fatih/color"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func createTestRecord() *TestLogRecord {
	return &TestLogRecord{
		TestStarted:        time.Now(),
		TxPerPayload:       viper.GetInt("env.txPerPayload"),
		TxIterations:       viper.GetInt("env.txIterations"),
		TxMultiBroadCast:   viper.GetInt("env.txMultiBroadCast"),
		TxDescription:      viper.GetString("env.txDescription"),
		CreatedAt:          time.Now(),
		ArkGoTesterVersion: ArkGoTesterVersion,
	}
}

func createTestLogTransaction(parentID int) *TestLogTransaction {
	return &TestLogTransaction{
		TestLogRecordID: parentID,
		CreatedAt:       time.Now(),
	}
}

func createTestIterationRecord(parentID int) *TestLogIteration {
	return &TestLogIteration{
		TestLogRecordID:  parentID,
		IterationStarted: time.Now(),
	}
}

//Save trx
func (tr *TestLogTransaction) Save() {
	err := ArkTestDB.Save(tr)
	if err != nil {
		log.Error("Error saving TestLogTransaction: ", err.Error())
	}
}

//Update or Save the record
func (tr *TestLogRecord) Update() {
	err := ArkTestDB.Update(tr)
	if err != nil {
		log.Error("Error updating TestLogRecord: ", err.Error())
	}
}

//Save the record
func (tr *TestLogRecord) Save() {
	err := ArkTestDB.Save(tr)
	if err != nil {
		log.Error("Error saving TestLogRecord: ", err.Error())
	}
}

//Update or Save the record
func (tr *TestLogIteration) Update() {
	err := ArkTestDB.Update(tr)
	if err != nil {
		log.Error("Error updating TestLogIteration: ", err.Error())
	}
}

//Save the record
func (tr *TestLogIteration) Save() {
	err := ArkTestDB.Save(tr)
	if err != nil {
		log.Error("Error saving TestLogIteration: ", err.Error())
	}
}

func listTestRecordsDB() {
	var results []TestLogRecord
	err := ArkTestDB.All(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, element := range results {
		fmt.Println(element)
	}
}

func listTestIterationsRecordsDB() {
	var results []TestLogIteration
	err := ArkTestDB.All(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, element := range results {
		fmt.Println(element)
	}
}

func getLatstTestRecord() (TestLogRecord, error) {
	var results []TestLogRecord
	err := ArkTestDB.All(&results, storm.Limit(1), storm.Reverse())

	rec := TestLogRecord{}

	if err != nil {
		log.Error(err.Error())
		return rec, err
	}
	rec = results[0]
	return rec, nil
}

func getTxIDsFromTestIterationRecords(testRec TestLogRecord) ([]string, error) {
	var results []TestLogIteration
	var err error
	var query storm.Query
	var transIDList []string

	query = ArkTestDB.Select(q.Eq("TestLogRecordID", testRec.ID))
	err = query.Find(&results)

	if err != nil {
		log.Error(err.Error())
		return transIDList, err
	}

	for _, el := range results {
		for _, txid := range el.TxIDs {
			transIDList = append(transIDList, txid)
		}
	}
	return transIDList, nil
}

func findConfirmations(testRec TestLogRecord) {
	transIDList, err := getTxIDsFromTestIterationRecords(testRec)
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

	for id, transIDPart := range divided {
		wg.Add(1)

		go func(transIDs []string, idPeer int, arkapi *core.ArkClient) {
			defer wg.Done()
			arkTmpClient := core.NewArkClientFromPeer(arkapi.GetRandomXPeers(1)[0])
			for _, txID := range transIDs {
				params := core.TransactionQueryParams{ID: txID}
				arkTransaction, _, _ := arkTmpClient.GetTransaction(params)

				confirmations := 0
				if arkTransaction.Success {
					confirmations = arkTransaction.SingleTransaction.Confirmations
				}

				rec := createTestLogTransaction(testRec.ID)
				rec.TransactionID = txID
				rec.Confirmations = confirmations
				rec.Save()
			}
		}(transIDPart, id, ArkAPIClient)
	}
}

func checkConfirmations(testRec TestLogRecord) {
	var results []TestLogTransaction
	var err error
	var query storm.Query

	fmt.Println("Calling REST endpoints to check confirmations....")
	wg.Wait()
	query = ArkTestDB.Select(q.Eq("TestLogRecordID", testRec.ID))
	err = query.Find(&results)

	if err != nil {
		log.Error(err.Error())
		return
	}

	fmt.Println("Checking confirmations...")
	missingCounter := 0
	for _, txRes := range results {
		if txRes.Confirmations < 1 {
			missingCounter++
			//fmt.Println(ix, ".ID=", txRes.TransactionID)
		}
	}
	if missingCounter == 0 {
		color.Set(color.FgGreen)
		fmt.Println("All sent transactions are confirmed...")
	} else {
		color.Set(color.FgRed)
		fmt.Println("Missing", missingCounter, " transactions. Try checking again in a minute or so...")
	}
}
