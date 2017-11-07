package core

import (
	"log"
	"testing"
)

func TestAutoConfigure(t *testing.T) {
	//SetActiveConfiguration(MAINNET)
	if len(EnvironmentParams.Network.Nethash) == 0 {
		t.Error("No NETWORK parameters read")
	}

	if EnvironmentParams.Fees.SecondSignature == 0 {
		t.Error("No FEES parameters read")
	}

	log.Println(t.Name(), EnvironmentParams.Network.Nethash)
	log.Println(t.Name(), EnvironmentParams.Fees.SecondSignature)
}

func TestGetConfigurationNative(t *testing.T) {
	arkapi := NewArkClient(nil)
	arkapi = arkapi.SetActiveConfiguration(MAINNET)
	log.Println(t.Name(), "Active network: ", EnvironmentParams.Network.Type, "BaseUrl", BaseURL)
	if EnvironmentParams.Network.Type != MAINNET {
		t.Error("Wrong network on init")
	}

	arkapi = arkapi.SetActiveConfiguration(DEVNET)
	log.Println(t.Name(), "Active network: ", EnvironmentParams.Network.Type, "BaseUrl", BaseURL)
	if EnvironmentParams.Network.Type != DEVNET {
		t.Error("Wrong network on init")
	}

	arkapi = arkapi.SetActiveConfiguration(MAINNET)
	log.Println(t.Name(), "Active network: ", EnvironmentParams.Network.Type, "BaseUrl", BaseURL)
	if EnvironmentParams.Network.Type != MAINNET {
		t.Error("Wrong network on init")
	}

	arkapi = arkapi.SetActiveConfiguration(KAPU)
	log.Println(t.Name(), "Active network: ", EnvironmentParams.Network.Type, "BaseUrl", BaseURL)
	if EnvironmentParams.Network.Type != KAPU {
		t.Error("Wrong network on init")
	}

}
