package core

import (
	"ark-go/arkcoin"
	"net/http"
	"strconv"

	"github.com/dghubble/sling"
)

var baseURL = ""

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
