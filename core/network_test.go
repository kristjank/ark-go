package core

import (
	"log"
	"os"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
}
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

func TestNewPeerArkApiClientFromIP(t *testing.T) {
	arkapi := NewArkClientFromIP("164.8.251.173")

	log.Println(arkapi.GetActivePeer())
	peerStatus, _, err := arkapi.GetConnectedPeerStatus()
	if peerStatus.Success {
		log.Println("New arkapi for peer", peerStatus)
		log.Println(t.Name(), "Success")
	} else {
		log.Println(t.Name(), "ERROR", err.Error())
	}
}
func TestGetRandomXPeers(t *testing.T) {
	arkapi := NewArkClient(nil)
	for i := 0; i < 20; i++ {
		arkapi = arkapi.SetActiveConfiguration(DEVNET)

		peers := arkapi.GetRandomXPeers(20)

		for _, el := range peers {
			log.Println(el)
		}
		arkapi = arkapi.SetActiveConfiguration(DEVNET)
	}
}
