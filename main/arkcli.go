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

	"github.com/fatih/color"
)

var arkclient = core.NewArkClient(nil)
var passphrase1 = ""
var passphrase2 = ""

func calculcateVoteRatio() core.TransactionPayload {
	arkapi := core.NewArkClient(nil)

	deleKey := "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		deleKey = "02bcfa0951a92e7876db1fb71996a853b57f996972ed059a950d910f7d541706c9 "
	}

	params := core.DelegateQueryParams{PublicKey: deleKey}

	votersEarnings := arkapi.CalculateVotersProfit(params, 0.70)

	var payload core.TransactionPayload

	//log.Println(t.Name(), "Success", votersEarnings)
	fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := float64(len(votersEarnings)) * (float64(core.EnvironmentParams.Fees.Send) / core.SATOSHI)
	for _, element := range votersEarnings {
		log.Println(fmt.Sprintf("|%s|%15.8f|%15.8f|%15.8f|%15.8f|%4d|%25d|",
			element.Address,
			element.VoteWeight,
			element.VoteWeightShare,
			element.EarnedAmount100,
			element.EarnedAmountXX,
			element.VoteDuration,
			int(element.EarnedAmountXX*core.SATOSHI)))

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		//transaction parameters
		tx := core.CreateTransaction(element.Address,
			int64(element.EarnedAmountXX*core.SATOSHI),
			"chris: 1st profit sharing payment... |tx made with ark-go",
			"",
			"")

		payload.Transactions = append(payload.Transactions, tx)

	}
	log.Println("Full forged amount: ", sumEarned, "Ratio calc check sum: ", sumRatio, "Amount to voters: ", sumShareEarned, "Ratio shared: ", float64(sumShareEarned)/float64(sumEarned), "Lottery:", int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI))
	log.Println(fmt.Sprintf("Payment fees: %2.2f", feeAmount))

	tx := core.CreateTransaction("ANqeL7CP2som7q9NFbRuaUc5WUnwYkSbFY",
		int64((sumEarned-sumShareEarned-feeAmount)*core.SATOSHI),
		"chris: 1st month lottery fund reserve... |tx made with ark-go",
		"",
		"")

	payload.Transactions = append(payload.Transactions, tx)

	return payload

	/*//payload complete - posting
	res, httpresponse, err := arkapi.PostTransaction(payload)
	if res.Success {
		log.Println("Success,", httpresponse.Status, res.TransactionIDs)

	} else {
		log.Println(res.Message, res.Error, httpresponse.Status, err.Error())

	}*/
}

func readAccountData() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter account passphrase: ")
	passphrase1, _ = reader.ReadString('\n')

	key := arkcoin.NewPrivateKeyFromPassword(passphrase1, arkcoin.ActiveCoinConfig)

	accountResp, _, _ := arkclient.GetAccount(core.AccountQueryParams{Address: key.PublicKey.Address()})
	if !accountResp.Success {
		log.Println("Error getting account data for address", key.PublicKey.Address())
		return
	}

	if accountResp.Account.SecondSignature == 1 {
		fmt.Print("Enter second account passphrase: ")
		passphrase2, _ = reader.ReadString('\n')
	}
}

//////////////////////////////////////////////////////////////////////////////

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls") //Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func printMenu() {
	CallClear()
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("banner1.txt")

	fmt.Print(string(dat))
	color.Unset()

	fmt.Println("==========================")
	fmt.Println("1-Display contributors")
	fmt.Println("2-Link account")
	fmt.Println("0-Exit")
	fmt.Println("==========================")
	color.Unset()
}

func main() {
	var choice int

	for choice != 9 {
		printMenu()
		fmt.Scan(&choice)
		switch choice {
		case 1:
			calculcateVoteRatio()
		case 2:

		case 3:

		case 4:
			fmt.Println("Bye!")
		default:
			fmt.Printf("Are you sure you have chosen properly?")
		}
	}
}
