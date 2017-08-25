package core

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	ErrorObj     error  `json:"ignore"`
}

//Error interface function
func (e ArkApiResponseError) Error() string {
	return fmt.Sprintf("ArkServiceApi: %v %v %v %v %v", e.ErrorObj.Error(), e.Success, e.ErrorMessage, e.Data, e.Message)
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

//NewArkClient creations with supported network
func NewArkClientForPeer(peer Peer) *ArkClient {
	BaseURL = "http://" + peer.IP + ":" + strconv.Itoa(peer.Port)
	EnvironmentParams.Network.ActivePeer = peer
	return &ArkClient{
		sling: sling.New().Base(BaseURL).
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

//SwitchPeer switches client connection to another node
//usage: must reassing pointer: arkapi = arkapi.SwitchPeer()
func (s *ArkClient) SwitchPeer() *ArkClient {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	//IF internal PeerList is empty - we do a full switch network - init from start
	if len(EnvironmentParams.Network.PeerList) == 0 {
		switchNetwork(EnvironmentParams.Network.Type)
		return NewArkClient(nil)
	}

	//if we have active memory peer list - we select a new random peer from already inited memlist
	//list is filled in LoadActiveConfiguration-where client init is made
	EnvironmentParams.Network.ActivePeer = EnvironmentParams.Network.PeerList[r1.Intn(len(EnvironmentParams.Network.PeerList))]
	BaseURL = "http://" + EnvironmentParams.Network.ActivePeer.IP + ":" + strconv.Itoa(EnvironmentParams.Network.ActivePeer.Port)

	//updating with latest peer data - setting height level
	resPeer, _, err := s.GetConnectedPeerStatus()
	if err == nil && resPeer.Success {
		EnvironmentParams.Network.ActivePeer.Height = resPeer.Header.Height
	}
	return NewArkClient(nil)
}

//GetActivePeer returns active peer connected
//doesn't call rest to update it
//updates when calling SwitchPeer or Network is Changed
func (s *ArkClient) GetActivePeer() Peer {
	return EnvironmentParams.Network.ActivePeer
}

//GetRandomXPeers returns requested number of randomly selected peers from healthy peer list
//used in multiplethreadbroadcast
func (s *ArkClient) GetRandomXPeers(numberOfPeers int) []Peer {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	nrAllPeers := len(EnvironmentParams.Network.PeerList)
	var retPeers []Peer
	for i := 0; i < numberOfPeers; i++ {
		retPeers = append(retPeers, EnvironmentParams.Network.PeerList[r1.Intn(nrAllPeers)])
	}
	return retPeers
}
