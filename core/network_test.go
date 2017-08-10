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

	arkapi.SwitchPeer()

	log.Println(t.Name(), "Success")
}
