package core

import (
	"ark-go/arkcoin"
	"net/http"

	"github.com/dghubble/sling"
)

var baseURL = ""

//ArkClient sling rest pointer
type ArkClient struct {
	sling *sling.Sling
}

func init() {
	baseURL = SetActiveConfiguration(DEVNET)
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
			Add("version", "1.0.1").
			Add("port", "4001").
			Add("Content-Type", "application/json"),
	}
}
