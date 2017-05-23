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
			Add("nethash", "6e84d08bd299ed97c212c886c98a57e36545c8f5d645ca7eeae63a8bd62d8988").
			Add("version", "1.0.1").
			Add("port", "4001").
			Add("Content-Type", "application/json"),
	}
}
