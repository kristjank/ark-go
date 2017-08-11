package core

import (
	"log"
	"os"
	"strconv"
	"testing"
)

func init() {
	log.SetOutput(os.Stdout)
}

func TestGetBlocks(t *testing.T) {
	arkapi := NewArkClient(nil)

	blockResponse, err, _ := arkapi.GetFullBlocksFromPeer(1512066)
	if blockResponse.Success {
		log.Println(t.Name(), "Success, returned ", len(blockResponse.Blocks), " blocks")
	} else {
		t.Error(err.Error())
	}
}

func TestGetPeerHeight(t *testing.T) {
	arkapi := NewArkClient(nil)

	blockResponse, err, _ := arkapi.GetPeerHeight()
	if blockResponse.Success {
		log.Println(t.Name(), "Success, returned block height", strconv.Itoa(blockResponse.Height), "ID", blockResponse.ID)
	} else {
		t.Error(err.Error())
	}
}
