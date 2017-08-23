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
