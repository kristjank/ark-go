package core

import (
	"net/http"
	"strconv"

	"github.com/kristjank/goark-node/api/model"
)

//GetFullBlocksFromPeer function returns a full list of blocks from current last block on. A random number of blocks is returned,
//due to ddos measures
func (s *ArkClient) GetFullBlocksFromPeer(lastBlockHeight int) (model.BlockResponse, ArkApiResponseError, *http.Response) {
	respData := new(model.BlockResponse)
	respError := new(ArkApiResponseError)

	resp, err := s.sling.New().Get("peer/blocks?lastBlockHeight="+strconv.Itoa(lastBlockHeight)).Receive(respData, respError)
	if err != nil {
		respError.ErrorMessage = err.Error()
	}

	return *respData, *respError, resp
}

//GetPeerHeight function returns node peer height.
func (s *ArkClient) GetPeerHeight() (model.BlockHeightResponse, ArkApiResponseError, *http.Response) {
	respError := new(ArkApiResponseError)
	respData := new(model.BlockHeightResponse)

	resp, err := s.sling.New().Get("api/blocks/getHeight").Receive(respData, respError)
	if err != nil {
		respError.ErrorMessage = err.Error()
	}

	return *respData, *respError, resp
}
