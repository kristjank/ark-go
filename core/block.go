package core

import (
	"net/http"
	"strconv"
)

//BlockResponse structure to receive blocks from a random peer
type BlockResponse struct {
	Success bool    `json:"success"`
	Blocks  []Block `json:"blocks"`
	Count   int     `json:"count"`
}

//BlockHeightResponse structure to receive blocks from a random peer
type BlockHeightResponse struct {
	Success bool   `json:"success"`
	Height  int    `json:"height"`
	ID      string `json:"id"`
}

//Block structure
type Block struct {
	ID                   string        `json:"id" storm:"id"`
	Version              int           `json:"version"`
	Timestamp            int           `json:"timestamp"`
	Height               int           `json:"height" storm:"index,unique"`
	PreviousBlock        string        `json:"previousBlock"`
	NumberOfTransactions int           `json:"numberOfTransactions"`
	TotalAmount          int           `json:"totalAmount"`
	TotalFee             int           `json:"totalFee"`
	Reward               int           `json:"reward"`
	PayloadLength        int           `json:"payloadLength"`
	PayloadHash          string        `json:"payloadHash"`
	GeneratorPublicKey   string        `json:"generatorPublicKey"`
	GeneratorID          string        `json:"generatorId"`
	BlockSignature       string        `json:"blockSignature"`
	Confirmations        int           `json:"confirmations"`
	TotalForged          string        `json:"totalForged"`
	Transactions         []Transaction `json:"transactions,omitempty"`
}

//GetFullBlocks function returns a full list of blocks. A radnom number of blocks is returned,
//due to ddos measures
func (s *ArkClient) GetFullBlocks(lastBlockHeight int) (BlockResponse, *http.Response, error) {
	respData := new(BlockResponse)
	respError := new(ArkApiResponseError)

	qstr := "lastBlockHeight=" + strconv.Itoa(lastBlockHeight)

	resp, err := s.sling.New().Get("api/blocks/?"+qstr).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}

//GetPeerHeight function returns node peer height.
func (s *ArkClient) GetPeerHeight() (BlockHeightResponse, *http.Response, error) {
	respError := new(ArkApiResponseError)
	respData := new(BlockHeightResponse)

	resp, err := s.sling.New().Get("api/blocks/getHeight").Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}
