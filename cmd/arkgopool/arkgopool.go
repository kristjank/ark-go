package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/kristjank/ark-go/arkcoin"
	"github.com/kristjank/ark-go/core"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/asdine/storm"
)

var arkclient = core.NewArkClient(nil)
var reader = bufio.NewReader(os.Stdin)
var arkpooldb *storm.DB
var wg sync.WaitGroup
var wgConfirmations sync.WaitGroup

var version = "master"

func initLogger() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("log/arkgo-pool.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(io.MultiWriter(file))
	} else {
		log.Error("Failed to log to file, using default stderr")
	}

}

func initializeBoltClient() {
	var err error
	arkpooldb, err = storm.Open(viper.GetString("client.dbfilename"))

	if err != nil {
		log.Info("Error in initializeBoldClient", err.Error())
		log.Panic(err.Error())
	}

	log.Println("DB Opened succefully")
	//defer arkpooldb.Close()
}

func readAccountData() (string, string) {
	fmt.Println("\nEnter account passphrase")
	fmt.Print("-->")
	pass1, _ := reader.ReadString('\n')
	re := regexp.MustCompile("\r?\n")
	pass1 = re.ReplaceAllString(pass1, "")

	pass2 := ""
	key := arkcoin.NewPrivateKeyFromPassword(pass1, arkcoin.ActiveCoinConfig)

	accountResp, _, _ := arkclient.GetAccount(core.AccountQueryParams{Address: key.PublicKey.Address()})
	deleResp, _, _ := arkclient.GetDelegate(core.DelegateQueryParams{PublicKey: string(key.PublicKey.Serialize())})
	if !accountResp.Success {
		log.Info("Error getting account data for delegate: " + deleResp.SingleDelegate.Username + "[" + key.PublicKey.Address() + "]")
		return "error", ""
	}

	if accountResp.Account.SecondSignature == 1 {
		fmt.Println("\nEnter second account passphrase for delegate: " + deleResp.SingleDelegate.Username + "[" + key.PublicKey.Address() + "]")
		fmt.Print("-->")
		pass2, _ = reader.ReadString('\n')
		re := regexp.MustCompile("\r?\n")
		pass2 = re.ReplaceAllString(pass2, "")
	}

	return pass1, pass2
}

func loadConfig(configFile string) {
	viper.SetConfigName(configFile) // name of config file (without extension)
	viper.AddConfigPath("settings") // path to look for the config file in
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file

	if err != nil {
		log.Info("No productive config found - loading sample")
		// try to load sample config
		viper.SetConfigName("sample.config")
		viper.AddConfigPath("settings")
		err := viper.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			log.Info("No configuration file loaded - using defaults")
		}
	}

	viper.SetDefault("delegate.address", "")
	viper.SetDefault("delegate.pubkey", "")

	viper.SetDefault("voters.shareRatio", 0.0)
	viper.SetDefault("voters.txdescription", "share tx by ark-go")
	viper.SetDefault("voters.fidelity", true)
	viper.SetDefault("voters.fidelityLimit", 24)
	viper.SetDefault("voters.minamount", 0.0)
	viper.SetDefault("voters.minVoteTime", 0)
	viper.SetDefault("voters.deductTxFees", true)
	viper.SetDefault("voters.blocklist", "")
	viper.SetDefault("voters.capBalance", false)
	viper.SetDefault("voters.balanceCapAmount", 0.0)
	viper.SetDefault("voters.whitelist", "")
	viper.SetDefault("voters.blockBalanceCap", true)

	viper.SetDefault("costs.address", "")
	viper.SetDefault("costs.shareRatio", 0.0)
	viper.SetDefault("costs.txdescription", "cost tx by ark-go")

	viper.SetDefault("reserve.address", "")
	viper.SetDefault("reserve.shareRatio", 0.0)
	viper.SetDefault("reserve.txdescription", "reserve tx by ark-go")

	viper.SetDefault("personal.address", "")
	viper.SetDefault("personal.shareRatio", 0.0)
	viper.SetDefault("personal.txdescription", "personal tx by ark-go")

	viper.SetDefault("client.network", "DEVNET")
	viper.SetDefault("client.dbFilename", "payment.db")
	viper.SetDefault("client.multibroadcast", 10)
	viper.SetDefault("client.payloadsize", 30)
	viper.SetDefault("client.autoconfigPeer", "")
	viper.SetDefault("client.statistics", false)
	viper.SetDefault("client.statPeer", "")
	viper.SetDefault("client.statPort", 54010)

}

//////////////////////////////////////////////////////////////////////////////
//GUI RELATED STUFF
func pause() {
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Print("Press 'ENTER' key to continue... ")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
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

	fmt.Println("Connected on", core.EnvironmentParams.Network.Token, "peer:", core.BaseURL, "| ARKGoPool version", version)
	log.Info("Connected on ", core.EnvironmentParams.Network.Token, " peer: ", core.BaseURL, "| ARKGoPool version", version)
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("settings/banner.txt")
	fmt.Print(string(dat))
}

func printMenu() {
	log.Info("--------- MAIN MENU ----------------")
	clearScreen()
	printBanner()
	printNetworkInfo()
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("\t1-Display contributors")
	fmt.Println("\t2-Send reward payments")
	fmt.Println("\t3-Link account")
	fmt.Println("\t4-List payment history")
	fmt.Println("\t5-Send bonus payments")
	fmt.Println("\t0-Exit")
	fmt.Println("")
	fmt.Print("\tSelect option [0-5]:")
	color.Unset()
}

func main() {
	// Reading input parameters
	configPtr := flag.String("config", "config", "Name of config file to use (without extension)")
	modeSilentPtr := flag.Bool("silent", false, "Is silent mode")
	flag.Parse()

	// Load configuration and defaults
	// Order is important
	loadConfig(*configPtr)
	initLogger()

	log.Info("=============================================================================")
	log.Info("ARKGO client starting")

	log.Info("Starting initialize Bolt client.")
	initializeBoltClient()
	log.Info("Bolt client initialized.")

	if len(viper.GetString("client.autoconfigPeer")) > 0 {
		log.Info("ARKGO client setting properties via autoconfig peer ", viper.GetString("client.autoconfigPeer"))
		arkclient = arkclient.SetActiveConfigurationFromPeerAddress(viper.GetString("client.autoconfigPeer"))
	} else if viper.GetString("client.network") == "DEVNET" {
		arkclient = arkclient.SetActiveConfiguration(core.DEVNET)
	} else if viper.GetString("client.network") == "KAPU" {
		arkclient = arkclient.SetActiveConfiguration(core.KAPU)
	} else {
		arkclient = arkclient.SetActiveConfiguration(core.MAINNET)
	}

	log.Info("ArkApiClient connected, active peer: ", arkclient.GetActivePeer())

	if *modeSilentPtr {
		log.Info("Silent Mode active")
		log.Info("Starting to send payments")
		SendPayments(true)
		log.Info("Waiting for threads to complete")
		color.Unset()
		wg.Wait()
		log.Info("Exiting silent mode and arkgopool")

		os.Exit(1985)
	}

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
			DisplayCalculatedVoteRatio()
			color.Unset()
		case 2:
			clearScreen()
			color.Set(color.FgHiGreen)
			SendPayments(false)
			wg.Wait()
			color.Unset()
		case 3:
			clearScreen()
			save(readAccountData())
			color.Set(color.FgHiGreen)
			log.Info("Account successfully linked")
			fmt.Println("Account successfully linked")
			pause()
			color.Unset()
		case 5:
			clearScreen()
			color.Set(color.FgHiGreen)

			fmt.Println("\nEnter bonus amount to send to loyal voters")
			fmt.Print("-->")
			sAmount2Send, err := reader.ReadString('\n')
			re := regexp.MustCompile("\r?\n")
			sAmount2Send = re.ReplaceAllString(sAmount2Send, "")

			fmt.Println("\nEnter bonus transaction description (vendor field)")
			fmt.Print("-->")

			txBonusDesc, err := reader.ReadString('\n')
			txBonusDesc = re.ReplaceAllString(txBonusDesc, "")

			iAmount2Send, err := strconv.Atoi(sAmount2Send)
			if err != nil {
				log.Error("Stopping bonus payment", err.Error())
				return
			}

			SendBonusPayment(iAmount2Send, txBonusDesc)
			pause()
			color.Unset()
		case 4:
			clearScreen()
			color.Set(color.FgHiGreen)
			listPaymentsDB()
			pause()
			//listPaymentsDetailsFromDB()
			//pause()
			color.Unset()
		}
	}
	color.Unset()
	log.Info("Exiting arkgopool....")
}
