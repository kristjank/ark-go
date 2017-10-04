package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/kristjank/ark-go/core"

	"github.com/fatih/color"
)

var arkclient = core.NewArkClient(nil)
var reader = bufio.NewReader(os.Stdin)

// ArkGoRepVersion v0.0.X
var ArkGoRepVersion string

func initLogger() {
	// Log as JSON instead of the default ASCII formatter.

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("log/arkgorep.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(file))
	} else {
		log.Error("Failed to log to file, using default stderr")
	}

}

func findHopper() {

	color.Set(color.FgHiGreen)
	fmt.Println("\nEnter your delegate address")
	fmt.Print("-->")
	delegateAddress, err := reader.ReadString('\n')
	re := regexp.MustCompile("\r?\n")
	delegateAddress = re.ReplaceAllString(delegateAddress, "")

	accountResp, _, _ := arkclient.GetAccount(core.AccountQueryParams{Address: delegateAddress})
	params := core.DelegateQueryParams{PublicKey: accountResp.Account.PublicKey}
	votersResp, _, _ := arkclient.GetDelegateVoters(params)

	for ix, element := range votersResp.Accounts {

		balance, _ := strconv.Atoi(element.Balance)

		//increase how many TX to check (gets slower)
		sendHistory := 3
		receiveHistory := 3

		sendTransResp, _, _ := arkclient.ListTransaction(core.TransactionQueryParams{SenderID: element.Address, Limit: sendHistory})
		receiveTransResp, _, _ := arkclient.ListTransaction(core.TransactionQueryParams{RecipientID: element.Address, Limit: receiveHistory})

		hopper := false
		lastSent := 0
		//Big movements out could be hoppers
		for idx, element := range sendTransResp.Transactions {
			if element.Amount > int64(balance)*9/10 {
				hopper = true
			}
			if idx == 0 {
				lastSent = int(element.Amount)
			}
		}

		// If balance is now null we should compare to what was moved out last
		if balance < 1 {
			balance = lastSent
		}

		//But not of most last transactions in where small
		for _, element := range receiveTransResp.Transactions {
			if element.Amount < int64(balance)*2/10 {
				hopper = false
			}
		}

		if hopper {
			logRow := fmt.Sprintf("%3d.|%s|%s", ix+1, element.Address, element.Balance)
			fmt.Println(logRow)
			log.Info(logRow)
		}

	}

	println("Finished finding hopper for address", delegateAddress, "check files in output folder")

	if err != nil {
		log.Error("Stopping find hopper", err.Error())
		return
	}
}

//////////////////////////////////////////////////////////////////////////////
//GUI RELATED STUFF
func pause() {
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Print("Press 'ENTER' key to continue... ")
	reader.ReadString('\n')
}

func clearScreen() {
	cmd := exec.Command("clear")
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()

}

func printNetworkInfo() {
	color.Set(color.FgHiCyan)
	if core.EnvironmentParams.Network.Type == core.MAINNET {
		fmt.Println("Connected to MAINNET on peer:", core.BaseURL, "| ArkGoRep version", ArkGoRepVersion)
	}

	if core.EnvironmentParams.Network.Type == core.DEVNET {
		fmt.Println("Connected to DEVNET on peer:", core.BaseURL, "| ArkGoRep version", ArkGoRepVersion)
	}
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("banner.txt")
	fmt.Print(string(dat))
}

func printMenu() {
	log.Info("--------- MAIN MENU ----------------")
	clearScreen()
	printBanner()
	printNetworkInfo()
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("\t1-Find Hopper")
	fmt.Println("\t2-Option 2")
	fmt.Println("\t3-Option 3")
	fmt.Println("\t4-Option 4")
	fmt.Println("\t0-Exit")
	fmt.Println("")
	fmt.Print("\tSelect option [1-4]:")
	color.Unset()
}

func main() {
	//sending ARKGO Server that we are working with payments
	//setting the version
	ArkGoRepVersion = "v0.0.1"

	// Load configration and defaults
	// Order is important
	initLogger()

	log.Info("Arkgorep client starting")
	log.Info("Arkgorep connected, active peer: ", arkclient.GetActivePeer())

	//switch to preset network
	arkclient = arkclient.SetActiveConfiguration(core.MAINNET)

	var choice = 1
	for choice != 0 {
		//pause()
		printMenu()

		//fmt.Scan(&choice)
		fmt.Fscan(reader, &choice)
		reader.ReadString('\n')

		switch choice {
		case 1:
			clearScreen()
			color.Set(color.FgMagenta)
			findHopper()
			pause()
			color.Unset()
		case 2:
			clearScreen()
			color.Set(color.FgHiGreen)
			//DO
			pause()
			color.Unset()
		case 3:
			clearScreen()
			color.Set(color.FgMagenta)
			//DO
			pause()
			color.Unset()
		case 4:
			clearScreen()
			color.Set(color.FgHiGreen)
			//DO
			pause()
			color.Unset()
		}
	}
	color.Unset()
}
