package core

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

const baseURL = "http://164.8.251.173:4001"

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

// ArkClientParams - when set, they are automatically added to get requests
type ArkClientParams struct {
	Address   string `url:"address,omitempty"`
	ID        string `url:"id,omitempty"`
	UserName  string `url:"username,omitempty"`
	PublicKey string `url:"publickey,omitempty"`
	IP        string `url:"ip,omitempty"`
	Port      string `url:"port,omitempty"`
}

//ArkClientError struct to hold error response
type ArkClientError struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"error"`
}

func (e ArkClientError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v", e.Success, e.ErrorMessage)
}

//ArkClient sling rest pointer
type ArkClient struct {
	sling *sling.Sling
}

//NewArkClient creations
func NewArkClient(httpClient *http.Client) *ArkClient {
	return &ArkClient{
		sling: sling.New().Client(httpClient).Base(baseURL).
			Add("nethash", "6e84d08bd299ed97c212c886c98a57e36545c8f5d645ca7eeae63a8bd62d8988").
			Add("version", "1").
			Add("port", "4001").
			Add("Content-Type", "application/json"),
	}
}

//ListPeers function returns list of peers from ArkNode
func (s *ArkClient) ListPeers(params *ArkClientParams) (PeerResponse, *http.Response, error) {
	peerResponse := new(PeerResponse)
	arkClientError := new(ArkClientError)
	resp, err := s.sling.New().Get("peer/list").QueryStruct(params).Receive(peerResponse, arkClientError)
	if err == nil {
		err = arkClientError
	}

	return *peerResponse, resp, err
}
