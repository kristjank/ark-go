package core

import (
	"log"
	"testing"
)

func TestListPeers(t *testing.T) {
	arkapi := NewArkClient(nil)
	peers, _, err := arkapi.ListPeers(nil)

	if peers.Success {
		log.Println(t.Name(), "Success, returned ", len(peers.Peers), "peers")
	} else {
		t.Error(err.Error())
	}
}
