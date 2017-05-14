package cmd

import (
	"log"
	"testing"
)

func TestKeys(t *testing.T) {
	address := GetAddress("this is a top secret passphrase")
	log.Println("address=", address)

}
