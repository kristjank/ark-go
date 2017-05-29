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

	params := PeerQueryParams{IP: "137.74.90.194", Port: 4001}

	peersResp, _, err := arkapi.GetPeer(params)
	if peersResp.Success {
		log.Println(t.Name(), "Success, returned ", peersResp.SinglePeer.Os, peersResp.SinglePeer.Status)

	} else {
		t.Error(err.Error())
	}
}
