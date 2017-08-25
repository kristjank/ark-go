package core

import (
	"log"
	"testing"
)

func TestNewArkClient(t *testing.T) {
	arkapi := NewArkClient(nil)

	if arkapi == nil {
		t.Error("Error creating client")
	}
	log.Println(t.Name(), "Success")
}

func TestSwitchPeer(t *testing.T) {
	arkapi := NewArkClient(nil)

	if arkapi == nil {
		t.Error("Error creating client")
	}

	arkapi = arkapi.SwitchPeer()
	log.Println(arkapi.GetActivePeer())
	log.Println(t.Name(), "Success")
}

func TestSwitchNetwork(t *testing.T) {
	arkapi := NewArkClient(nil)

	for i := 0; i < 5; i++ {
		arkapi = arkapi.SetActiveConfiguration(MAINNET)
		log.Println("Selected: ", arkapi.GetActivePeer())
	}

	arkapi = arkapi.SwitchPeer()
	log.Println(arkapi.GetActivePeer())
	log.Println(t.Name(), "Success")
}

func TestNewPeerArkApiClient(t *testing.T) {
	arkapi := NewArkClient(nil)

	log.Println(arkapi.GetActivePeer())
	peer := EnvironmentParams.Network.PeerList[3]
	arkapi1 := NewArkClientFromPeer(peer)
	log.Println("New arkapi for peer", arkapi1.GetActivePeer())

	log.Println(arkapi1.GetActivePeer())
	log.Println(t.Name(), "Success")
}

func TestGetRandomXPeers(t *testing.T) {
	arkapi := NewArkClient(nil)
	for i := 0; i < 20; i++ {
		arkapi = arkapi.SetActiveConfiguration(DEVNET)

		peers := arkapi.GetRandomXPeers(5)

		for _, el := range peers {
			log.Println(el)
		}
		arkapi = arkapi.SetActiveConfiguration(DEVNET)
	}
}
