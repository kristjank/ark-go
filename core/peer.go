package core

import (
	"fmt"
	"net/http"
)

//PeerResponse structure for call /peer/list
type PeerResponse struct {
	Success bool `json:"success"`
	Peers   []struct {
		IP      string `json:"ip"`
		Port    int    `json:"port"`
		Version string `json:"version,omitempty"`
		Os      string `json:"os,omitempty"`
		Height  int    `json:"height,omitempty"`
		Status  string `json:"status"`
		Delay   int    `json:"delay"`
	} `json:"peers"`
}

// PeerQueryParams - when set, they are automatically added to get requests
type PeerQueryParams struct {
	State   string `url:"state,omitempty"`   //State of peer. 1 - disconnected. 2 - connected. 0 - banned. (String)
	Os      string `url:"os,omitempty"`      //OS of peer. (String)
	Shared  string `url:"shared,omitempty"`  //Is peer shared? Boolean: true or false. (String)
	Version string `url:"version,omitempty"` //Version of peer. (String)
	Limit   string `url:"limit,omitempty"`   //Limit to show. Max limit is 100. (Integer)
	OrderBy string `url:"orderBy,omitempty"` //Name of column to order. After column name must go 'desc' or 'asc' to choose order type. (String)
}

//PeerResponseError struct to hold error response
type PeerResponseError struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

//Error interface function
func (e PeerResponseError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v", e.Success, e.ErrorMessage)
}

//ListPeers function returns list of peers from ArkNode
func (s *ArkClient) ListPeers(params *PeerQueryParams) (PeerResponse, *http.Response, error) {
	peerResponse := new(PeerResponse)
	peerResponseError := new(PeerResponseError)
	resp, err := s.sling.New().Get("peer/list").QueryStruct(params).Receive(peerResponse, peerResponseError)
	if err == nil {
		err = peerResponseError
	}

	return *peerResponse, resp, err
}
