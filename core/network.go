package core

import (
	"ark-go/arkcoin"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dghubble/sling"
)

var baseURL = ""

//ArkApiResponseError struct to hold error response from api node
type ArkApiResponseError struct {
	Success      bool   `json:"success,omitempty"`
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
	Data         string `json:"data,omitempty"`
}

//Error interface function
func (e ArkApiResponseError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v", e.Success, e.ErrorMessage, e.Data, e.Message)
}

//ArkClient sling rest pointer
type ArkClient struct {
	sling *sling.Sling
}

func init() {
	baseURL = SetActiveConfiguration(MAINNET)
	coinParams := arkcoin.Params{
		AddressHeader: EnvironmentParams.Network.AddressVersion,
	}
	arkcoin.SetActiveCoinConfiguration(&coinParams)
}

//NewArkClient creations
func NewArkClient(httpClient *http.Client) *ArkClient {

	return &ArkClient{
		sling: sling.New().Client(httpClient).Base(baseURL).
			Add("nethash", EnvironmentParams.Network.Nethash).
			Add("version", EnvironmentParams.Network.ActivePeer.Version).
			Add("port", strconv.Itoa(EnvironmentParams.Network.ActivePeer.Port)).
			Add("Content-Type", "application/json"),
	}
}
