package core

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/sling"
)

//BaseURL holds the IP,PORT of the connected peer used for the communication
var BaseURL = ""

//ArkApiResponseError struct to hold error response from api node
type ArkApiResponseError struct {
	Success      bool   `json:"success,omitempty"`
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
	Data         string `json:"data,omitempty"`
}

//Error interface function
func (e ArkApiResponseError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v %v %v", e.Success, e.ErrorMessage, e.Data, e.Message)
}

//ArkClient sling rest pointer
type ArkClient struct {
	sling *sling.Sling
}

func init() {
	switchNetwork(MAINNET)
}

//NewArkClient creations with supported network
func NewArkClient(httpClient *http.Client) *ArkClient {
	return &ArkClient{
		sling: sling.New().Client(httpClient).Base(BaseURL).
			Add("nethash", EnvironmentParams.Network.Nethash).
			Add("version", EnvironmentParams.Network.ActivePeer.Version).
			Add("port", strconv.Itoa(EnvironmentParams.Network.ActivePeer.Port)).
			Add("Content-Type", "application/json"),
	}
}

//TestMethodNewArkClient creations with supported network
//A test method for local node testing when implementid
//Not for production use
func TestMethodNewArkClient(httpClient *http.Client) *ArkClient {
	return &ArkClient{
		sling: sling.New().Client(httpClient).Base("http://164.8.251.173:4001").
			Add("nethash", EnvironmentParams.Network.Nethash).
			Add("version", EnvironmentParams.Network.ActivePeer.Version).
			Add("port", strconv.Itoa(EnvironmentParams.Network.ActivePeer.Port)).
			Add("Content-Type", "application/json"),
	}
}
