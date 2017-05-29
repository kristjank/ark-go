package cmd

import (
	"log"
	"testing"
)

func TestKeys(t *testing.T) {
	address := GetAddress("this is chrises top secret dev account passphrase")

	log.Println("address=", address)

}
