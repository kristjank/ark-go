package core

import "testing"

func TestAutoConfigure(t *testing.T) {
	arkapi := NewArkClient(nil)

	arkapi.AutoConfigureParams()
}
