package core

import (
	"net/http"
	"strconv"

	"github.com/kristjank/goark-node/api/model"
)

//GetFullBlocksFromPeer function returns a full list of blocks. A radnom number of blocks is returned,
//due to ddos measures
func (s *ArkClient) GetFullBlocksFromPeer(lastBlockHeight int) (model.BlockResponse, *http.Response, error) {
	respData := new(model.BlockResponse)
	respError := new(ArkApiResponseError)

	qstr := "lastBlockHeight=" + strconv.Itoa(lastBlockHeight)

	resp, err := s.sling.New().Get("api/blocks/?"+qstr).Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}

//GetPeerHeight function returns node peer height.
func (s *ArkClient) GetPeerHeight() (model.BlockHeightResponse, *http.Response, error) {
	respError := new(ArkApiResponseError)
	respData := new(model.BlockHeightResponse)

	resp, err := s.sling.New().Get("api/blocks/getHeight").Receive(respData, respError)
	if err == nil {
		err = respError
	}

	return *respData, resp, err
}
