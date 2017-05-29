package core

import (
	"log"
	"testing"
)

func TestAutoConfigure(t *testing.T) {
	SetActiveConfiguration(MAINNET)
	if len(EnvironmentParams.Network.Nethash) == 0 {
		t.Error("No NETWORK parameters read")
	}

	if EnvironmentParams.Fees.Secondsignature == 0 {
		t.Error("No FEES parameters read")
	}

	log.Println(t.Name(), EnvironmentParams.Network.Nethash)
	log.Println(t.Name(), EnvironmentParams.Fees.Secondsignature)
}

func TestGetConfigurationNative(t *testing.T) {
	//GetConfigurationNative()
}
