package core

import (
	"log"
	"testing"
)

func TestListPeers(t *testing.T) {
	arkapi := NewArkClient(nil)

	params := PeerQueryParams{Status: "OK", Version: "1.0.1"}

	peers, _, err := arkapi.ListPeers(params)
	if peers.Success {
		log.Println(t.Name(), "Success, returned ", len(peers.Peers), "peers")

	} else {
		t.Error(err.Error())
	}
}

func TestGetPeer(t *testing.T) {
	arkapi := NewArkClient(nil)
	ipAddress := "137.74.90.194"
	portNum := 4001

	if EnvironmentParams.Network.Type == DEVNET {
		ipAddress = "164.8.251.91"
		portNum = 4002
	}

	params := PeerQueryParams{IP: ipAddress, Port: portNum}

	peersResp, _, err := arkapi.GetPeer(params)
	if peersResp.Success {
		log.Println(t.Name(), "Success, returned ", peersResp.SinglePeer.Os, peersResp.SinglePeer.Status)

	} else {
		t.Error(err.Error())
	}
}
