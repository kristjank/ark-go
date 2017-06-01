package main

import (
	"ark-go/arkcoin"
	"ark-go/core"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/viper"
)

var arkclient = core.NewArkClient(nil)

var reader = bufio.NewReader(os.Stdin)

var errorlog *os.File
var logger *log.Logger

func init() {
	errorlog, err := os.OpenFile("ark-goclient.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.New(errorlog, "applog: ", log.Lshortfile|log.LstdFlags)
}

//DisplayCalculatedVoteRatio based on parameters in config.toml
func DisplayCalculatedVoteRatio() {
	pubKey := viper.GetString("delegate.pubkey")
	params := core.DelegateQueryParams{PublicKey: pubKey}
	deleResp, _, _ := arkclient.GetDelegate(params)
	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"))
	shareRatioStr := strconv.FormatFloat(viper.GetFloat64("voters.shareratio")*100, 'f', -1, 64) + "%"

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := float64(len(votersEarnings)) * (float64(core.EnvironmentParams.Fees.Send) / core.SATOSHI)
	color.Set(color.FgHiGreen)
	fmt.Println("Displaying information for delegate:", deleResp.SingleDelegate.Username, deleResp.SingleDelegate.Address)
	fmt.Println("")
	fmt.Println(fmt.Sprintf("|%34s|%18s|%8s|%17s|%17s|%6s|", "Address", "Balance", "Weight", "Reward-100%", "Reward-"+shareRatioStr, "Hours"))
	color.Set(color.FgCyan)
	for _, element := range votersEarnings {
		s := fmt.Sprintf("|%s|%18.8f|%8.4f|%15.8f A|%15.8f A|%6d|",
			element.Address,
			element.VoteWeight,
			element.VoteWeightShare,
			element.EarnedAmount100,
			element.EarnedAmountXX,
			element.VoteDuration)

		fmt.Println(s)
		logger.Println(s)

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}

	fmt.Println("")
	fmt.Println("..............")
	fmt.Println("Available amount:", sumEarned)
	fmt.Println("Ratio calc check:", sumRatio)
	fmt.Println("Amount to voters:", sumShareEarned)
	fmt.Println("Ratio share check:", float64(sumShareEarned)/float64(sumEarned))
	fmt.Println("Lottery:", int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI))
	fmt.Println(fmt.Sprintf("Payment fees: %2.2f", feeAmount))

	logger.Println("Full forged amount: ", sumEarned, "Ratio calc check sum: ", sumRatio, "Amount to voters: ", sumShareEarned, "Ratio shared: ", float64(sumShareEarned)/float64(sumEarned), "Lottery:", int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI))
	logger.Println(fmt.Sprintf("Payment fees: %2.2f", feeAmount))

	pause()
}

//SendPayments based on parameters in config.toml
func SendPayments() {
	pubKey := viper.GetString("delegate.pubkey")
	params := core.DelegateQueryParams{PublicKey: pubKey}
	var payload core.TransactionPayload

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"))

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := float64(len(votersEarnings)+1) * (float64(core.EnvironmentParams.Fees.Send) / core.SATOSHI)

	p1, p2 := readAccountData()

	for _, element := range votersEarnings {
		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		//transaction parameters
		tx := core.CreateTransaction(element.Address,
			int64(element.EarnedAmountXX*core.SATOSHI),
			viper.GetString("voters.txdescription"),
			p1,
			p2)

		payload.Transactions = append(payload.Transactions, tx)
	}

	tx := core.CreateTransaction("ANqeL7CP2som7q9NFbRuaUc5WUnwYkSbFY",
		int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI),
		"chris: 1st month lottery fund reserve... |tx made with ark-go",
		p1,
		p2)

	payload.Transactions = append(payload.Transactions, tx)

	pause()

	/*//payload complete - posting
	res, httpresponse, err := arkapi.PostTransaction(payload)
	if res.Success {
		log.Println("Success,", httpresponse.Status, res.TransactionIDs)

	} else {
		log.Println(res.Message, res.Error, httpresponse.Status, err.Error())

	}*/
}

func readAccountData() (string, string) {
	reader.ReadString('\n')
	fmt.Print("Enter account passphrase: ")
	passphrase1, _ := reader.ReadString('\n')
	passphrase2 := ""

	key := arkcoin.NewPrivateKeyFromPassword(passphrase1, arkcoin.ActiveCoinConfig)

	accountResp, _, _ := arkclient.GetAccount(core.AccountQueryParams{Address: key.PublicKey.Address()})
	if !accountResp.Success {
		logger.Println("Error getting account data for address", key.PublicKey.Address())
		return "error", ""
	}

	if accountResp.Account.SecondSignature == 1 {
		reader.ReadString('\n')
		fmt.Print("Enter second account passphrase: ")
		passphrase2, _ = reader.ReadString('\n')
	}

	return passphrase1, passphrase2
}

//////////////////////////////////////////////////////////////////////////////
//GUI RELATED STUFF
func pause() {
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Print("Press any key to return to the menu... ")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	reader.ReadString('\n')
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
		fmt.Println("Connected to MAINNET on peer:", core.BaseURL)
	}

	if core.EnvironmentParams.Network.Type == core.DEVNET {
		fmt.Println("Connected to DEVNET on peer:", core.BaseURL)
	}
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("settings/banner5.txt")
	fmt.Print(string(dat))
}

func printMenu() {
	clearScreen()
	printBanner()
	printNetworkInfo()
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("\t1-Display contributors")
	fmt.Println("\t2-Send payments")
	fmt.Println("\t4-Link account")
	fmt.Println("\t8-Switch network")
	fmt.Println("\t9-Exit")
	fmt.Println("")
	fmt.Print("\tSelect option [1-9]:")
	color.Unset()
}

type cost struct {
	Address       string
	AddressRatio  float64
	TxDescription string
}

type costs struct {
	Cost []cost
}

func main() {
	logger.Println("Ark-golang client starting")

	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("settings") // path to look for the config file in
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	/*b := viper.AllSettings()
	fmt.Println(b)
	aa := viper.Get("costs")

	for key, value := range aa {

	}
	fmt.Println(aa)
	pause()*/

	var choice int

	for choice != 9 {
		//pause()
		printMenu()
		fmt.Scan(&choice)
		switch choice {
		case 1:
			clearScreen()
			color.Set(color.FgMagenta)
			DisplayCalculatedVoteRatio()
			color.Unset()
		case 2:
			clearScreen()
			color.Set(color.FgRed)
			SendPayments()
			color.Unset()

		case 3:
		case 8:
			if core.EnvironmentParams.Network.Type == core.MAINNET {
				arkclient = arkclient.SetActiveConfiguration(core.DEVNET)
			} else {
				arkclient = arkclient.SetActiveConfiguration(core.MAINNET)
			}
		}
	}

	defer errorlog.Close()
}
