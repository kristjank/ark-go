package core

import (
	"net/http"

	"github.com/dghubble/sling"
)

const baseURL = "http://5.39.9.240:4001"

//ArkClient sling rest pointer
type ArkClient struct {
	sling *sling.Sling
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
