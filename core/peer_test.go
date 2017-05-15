package core

import (
	"log"
	"testing"
)

func TestListPeers(t *testing.T) {
	arkapi := NewArkClient(nil)

	peers, e1, err := arkapi.ListPeers()

	if peers.Success {
		log.Println("Success")

		for _, v := range peers.Peers {
			log.Println(v.IP, v.Port)
		}

	} else {
		log.Println(e1.Status)
		log.Println(err.Error())
	}

	//log.Println(t, peers.Peers[5])

}
