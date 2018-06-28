package core

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/sling"
	"github.com/kristjank/ark-go/arkcoin"
)

//ArkNetworkType for network switching
type ArkNetworkType int

const (
	//MAINNET connection
	MAINNET = iota
	//DEVNET connection
	DEVNET
	//KAPU network
	KAPU
)

//TO HELP DIVIDE
//TODO rename
const (
	SATOSHI = 100000000
)

//EnvironmentParams - Global ARK EnvironmentParams read from acitve peers
var EnvironmentParams = new(ArkEnvParams)

//ArkEnvParams structure to hold parameters from autoconfigure
//structure is filled from peers at arkclient init call
type ArkEnvParams struct {
	Success bool    `json:"success"`
	Network Network `json:"network"`
	Fees    Fees    `json:"fees"`
}

//Fees constant parameters for active configuration
type Fees struct {
	Send            int64 `json:"send"`
	Vote            int64 `json:"vote"`
	SecondSignature int64 `json:"secondsignature"`
	Delegate        int64 `json:"delegate"`
	MultiSignature  int64 `json:"multisignature"`
}

//Network parameters
type Network struct {
	Nethash        string         `json:"nethash"`
	Token          string         `json:"token"`
	Symbol         string         `json:"symbol"`
	Explorer       string         `json:"explorer"`
	AddressVersion byte           `json:"version"` //this is address generator version!!!
	Type           ArkNetworkType //holding ark networktype
	ActivePeer     Peer
	PeerList       []Peer
}

//LoadActiveConfiguration reads arknetwork parameters from the Network
//and fills the EnvironmentParams structure
//selected and connected peer address is returned
func LoadActiveConfiguration(arknetwork ArkNetworkType) string {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	selectedPeer := ""
	EnvironmentParams.Network.Type = arknetwork
	//looping peers comunication until we get autoconfigure response
	switch arknetwork {
	case MAINNET:
		log.Println("Active network is ARK MAINNET")
		selectedPeer = seedList[r1.Intn(len(seedList))]
		log.Println("Random peer selected: ", selectedPeer)

	case DEVNET:
		log.Println("Active network is ARK DEVNET")
		selectedPeer = testSeedList[r1.Intn(len(testSeedList))]
		log.Println("Random peer selected: ", selectedPeer)
	case KAPU:
		log.Println("Active network is KAPU MAINNET")
		selectedPeer = seedListKAPU[r1.Intn(len(seedListKAPU))]
		log.Println("Random peer selected: ", selectedPeer)
	}

	return autoConfigFromPeer(selectedPeer)
}

func autoConfigFromPeer(seedPeer string) string {
	log.Println("Connecting to selected peer", seedPeer)
	//reading basic network params
	res, err := http.Get("http://" + seedPeer + "/api/loader/autoconfigure")
	if err != nil {
		log.Println("Error receiving autoloader params rest from: ", seedPeer, " Error: ", err.Error())
		seedPeer = ""
	} else {
		json.NewDecoder(res.Body).Decode(&EnvironmentParams)
	}

	if seedPeer == "" {
		log.Fatal("Unable to connect to blockchain, exiting")
	}

	//reading fees
	res, err = http.Get("http://" + seedPeer + "/api/blocks/getFees")
	if err != nil {
		log.Fatal("Error receiving fees params rest from: ", seedPeer)
	}
	json.NewDecoder(res.Body).Decode(&EnvironmentParams)

	//getting connected peer params from peer
	peerParams := strings.Split(seedPeer, ":")
	peerRes := new(PeerResponse)
	res, err = http.Get("http://" + seedPeer + "/api/peers/get/?ip=" + peerParams[0] + "&port=" + peerParams[1])
	if err != nil {
		log.Fatal("Error receiving peer status from: ", seedPeer, err.Error(), res.StatusCode)
	}
	json.NewDecoder(res.Body).Decode(peerRes)
	//saving peer parameters to globals
	EnvironmentParams.Network.ActivePeer = peerRes.SinglePeer

	coinParams := arkcoin.Params{
		AddressHeader: EnvironmentParams.Network.AddressVersion,
	}
	arkcoin.SetActiveCoinConfiguration(&coinParams)

	return "http://" + optimizePeerList(seedPeer)
}

func optimizePeerList(selectedPeer string) string {
	tmpClient := &ArkClient{
		sling: sling.New().Client(nil).Base("http://"+selectedPeer).
			Add("nethash", EnvironmentParams.Network.Nethash).
			Add("version", EnvironmentParams.Network.ActivePeer.Version).
			Add("port", strconv.Itoa(EnvironmentParams.Network.ActivePeer.Port)).
			Add("Content-Type", "application/json"),
	}

	peerResp, err, _ := tmpClient.GetAllPeers()
	if err.ErrorObj != nil {
		log.Println("Error getting peer list")
		return selectedPeer
	}

	EnvironmentParams.Network.PeerList = peerResp.Peers
	log.Println("Start to optimize peer list, currently ", len(EnvironmentParams.Network.PeerList), " peers.")

	//Clean the peer list (filters not working as they shoud) - so checking again here
	maxHeight := EnvironmentParams.Network.ActivePeer.Height
	for i := len(EnvironmentParams.Network.PeerList) - 1; i >= 0; i-- {
		peer := EnvironmentParams.Network.PeerList[i]

		// Condition to decide if current element has to be deleted:
		if peer.Status != "OK" {
			EnvironmentParams.Network.PeerList = append(EnvironmentParams.Network.PeerList[:i], EnvironmentParams.Network.PeerList[i+1:]...)
			//log.Println("Removing peer", peer.IP, peer.Status, peer.Height)
			continue
		}
		//if all is ok and height is higher - we preffer peers with higher hight
		if peer.Height > maxHeight {
			log.Println("Setting new active peer, found OK peer with bigger block height", peer.Height, maxHeight)
			EnvironmentParams.Network.ActivePeer = peer
			selectedPeer = fmt.Sprintf("%s:%d", peer.IP, peer.Port)
			maxHeight = peer.Height
		}
	}

	//removing peers with difference more then 10 blocks, that is 10x8s behing mainheight
	for i := len(EnvironmentParams.Network.PeerList) - 1; i >= 0; i-- {
		peer := EnvironmentParams.Network.PeerList[i]

		if maxHeight-peer.Height > 20 {
			EnvironmentParams.Network.PeerList = append(EnvironmentParams.Network.PeerList[:i], EnvironmentParams.Network.PeerList[i+1:]...)
			//log.Println("Removing peer, based on maxheight difference condition", peer.IP, peer.Status, peer.Height)
			continue
		}
	}
	log.Println("End of peer optimization, remaining ", len(EnvironmentParams.Network.PeerList), " peers.")
	return selectedPeer
}

func switchNetwork(arkNetwork ArkNetworkType) {
	BaseURL = LoadActiveConfiguration(arkNetwork)
	var wifHeader = []byte{170}

	coinParams := arkcoin.Params{
		AddressHeader:          EnvironmentParams.Network.AddressVersion,
		DumpedPrivateKeyHeader: wifHeader,
	}
	arkcoin.SetActiveCoinConfiguration(&coinParams)
}

//SetActiveConfiguration sets a new client connection, switches network and reads network settings from peer
//usage - must reassing new pointer value: arkapi = arkapi.SetActiveConfiguration(MAINNET)
func (s *ArkClient) SetActiveConfiguration(arkNetwork ArkNetworkType) *ArkClient {
	switchNetwork(arkNetwork)
	return NewArkClient(nil)
}

//SetActiveConfigurationFromPeerAddress sets a new client connection and autoconfigures from specified seed peer
func (s *ArkClient) SetActiveConfigurationFromPeerAddress(seedPeer string) *ArkClient {
	BaseURL = autoConfigFromPeer(seedPeer)
	return NewArkClient(nil)
}

var seedList = [...]string{
	"5.39.9.240:4001",
	"5.39.9.241:4001",
	"5.39.9.242:4001",
	"5.39.9.243:4001",
	"5.39.9.244:4001",
	"5.39.9.250:4001",
	"5.39.9.251:4001",
	"5.39.9.252:4001",
	"5.39.9.253:4001",
	"5.39.9.254:4001",
	"5.39.9.255:4001",
	"54.38.48.160:4001",
	"54.38.48.161:4001",
	"54.38.48.162:4001",
	"54.38.48.163:4001",
	"54.38.48.164:4001",
	"54.38.48.165:4001",
	"54.38.48.166:4001",
	"54.38.48.167:4001",
	"54.38.48.168:4001",
	"54.38.48.169:4001"}

var testSeedList = [...]string{
	"167.114.29.32:4002",
	"167.114.29.33:4002",
	"167.114.29.34:4002",
	"167.114.29.35:4002",
	"167.114.29.36:4002",
	"167.114.29.37:4002",
	"167.114.29.38:4002",
	"167.114.29.39:4002",
	"167.114.29.40:4002",
	"167.114.29.41:4002",
	"167.114.29.42:4002",
	"167.114.29.51:4002",
	"167.114.29.52:4002",
	"167.114.29.53:4002",
	"167.114.29.54:4002",
	"167.114.29.55:4002"}

var seedListKAPU = [...]string{
	"51.15.198.173:4600",
	"51.15.215.113:4600",
	"51.15.221.100:4600",
	"51.15.194.207:4600",
	"94.176.238.173:4600",
	"185.5.55.249:4600",
	"94.176.233.213:4600",
	"94.176.236.51:4600",
	"94.176.233.210:4600",
	"51.15.89.225:9700",
	"51.15.201.56:4600",
	"51.15.84.234:4600",
	"51.15.84.139:4600",
	"185.5.55.249:4600",
	"80.241.218.21:4600",
	"145.239.90.228:4600",
	"144.217.242.167:4600",
	"51.15.201.243:4600",
	"45.32.29.180:4600"}
