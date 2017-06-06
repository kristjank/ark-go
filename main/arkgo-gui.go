package main

import (
	"ark-go/arkcoin"
	"ark-go/core"
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
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
	errorlog, err := os.OpenFile("arkgo-gui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.New(errorlog, "ark-go: ", log.Lshortfile|log.LstdFlags)
}

//DisplayCalculatedVoteRatio based on parameters in config.toml
func DisplayCalculatedVoteRatio() {
	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}
	deleResp, _, _ := arkclient.GetDelegate(params)
	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"))
	shareRatioStr := strconv.FormatFloat(viper.GetFloat64("voters.shareratio")*100, 'f', -1, 64) + "%"

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("Displaying voter information for delegate:", deleResp.SingleDelegate.Username, deleResp.SingleDelegate.Address)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("|%-34s|%18s|%8s|%17s|%17s|%6s|", "Voter address", "Balance", "Weight", "Reward-100%", "Reward-"+shareRatioStr, "Hours"))
	color.Set(color.FgCyan)
	for _, element := range votersEarnings {
		s := fmt.Sprintf("|%s|%18.8f|%8.4f|%15.8f A|%15.8f A|%6d|", element.Address, element.VoteWeight, element.VoteWeightShare, element.EarnedAmount100, element.EarnedAmountXX, element.VoteDuration)

		fmt.Println(s)
		logger.Println(s)

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}

	//Cost calculation
	costAmount := sumEarned * viper.GetFloat64("costs.shareratio")
	reserveAmount := sumEarned * viper.GetFloat64("reserve.shareratio")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("Available amount:", sumEarned)
	fmt.Println("Amount to voters:", sumShareEarned, viper.GetFloat64("voters.shareratio"))
	fmt.Println("Amount to costs:", costAmount, viper.GetFloat64("costs.shareratio"))
	fmt.Println("Amount to reserve:", reserveAmount, viper.GetFloat64("reserve.shareratio"))

	fmt.Println("Ratio calc check:", sumRatio, "(should be = 1)")
	fmt.Println("Ratio share check:", float64(sumShareEarned)/float64(sumEarned), "should be=", viper.GetFloat64("voters.shareratio"))

	pause()
}

func checkConfigSharingRatio() bool {
	a1 := viper.GetFloat64("voters.shareratio")
	a2 := viper.GetFloat64("costs.shareratio")
	a3 := viper.GetFloat64("reserve.shareratio")

	if a1+a2+a3 != 1.0 {
		logger.Println("Wrong config. Check share ration percentages!")
		return false
	}
	return true
}

//SendPayments based on parameters in config.toml
func SendPayments() {
	if !checkConfigSharingRatio() {
		logger.Fatal("Unable to calculcate.")
	}

	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}
	var payload core.TransactionPayload

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"))

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0

	p1, p2 := readAccountData()
	clearScreen()

	for _, element := range votersEarnings {
		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		//transaction parameters
		txAmount2Send := int64(element.EarnedAmountXX*core.SATOSHI) - core.EnvironmentParams.Fees.Send
		tx := core.CreateTransaction(element.Address, txAmount2Send, viper.GetString("voters.txdescription"), p1, p2)

		payload.Transactions = append(payload.Transactions, tx)
	}

	//Cost & reserve fund calculation
	costAmount := sumEarned * viper.GetFloat64("costs.shareratio")
	reserveAmount := sumEarned * viper.GetFloat64("reserve.shareratio")

	//summary and conversion checks
	if (costAmount + reserveAmount + sumShareEarned) != sumEarned {
		color.Set(color.FgHiRed)
		log.Println("Calculation of ratios NOT OK - overall summary failing")
		logger.Println("Calculation of ratios NOT OK - overall summary failing")
	}

	costAmount2Send := int64(costAmount*core.SATOSHI) - core.EnvironmentParams.Fees.Send
	costAddress := viper.GetString("costs.address")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		costAddress = viper.GetString("costs.Daddress")
	}
	txCosts := core.CreateTransaction(costAddress, costAmount2Send, viper.GetString("costs.txdescription"), p1, p2)
	payload.Transactions = append(payload.Transactions, txCosts)

	reserveAddress := viper.GetString("reserve.address")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		reserveAddress = viper.GetString("reserve.Daddress")
	}
	reserveAmount2Send := int64(reserveAmount*core.SATOSHI) - core.EnvironmentParams.Fees.Send

	txReserve := core.CreateTransaction(reserveAddress, reserveAmount2Send, viper.GetString("reserve.txdescription"), p1, p2)
	payload.Transactions = append(payload.Transactions, txReserve)

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("Transactions to be sent:")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	color.Set(color.FgHiCyan)
	for _, el := range payload.Transactions {
		s := fmt.Sprintf("|%s|%15d| %-40s|", el.RecipientID, el.Amount, el.VendorField)
		fmt.Println(s)
		logger.Println(s)
	}

	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Print("Send transactions and complete reward payments [Y/N]: ")

	c, _ := reader.ReadByte()

	if c == []byte("Y")[0] || c == []byte("y")[0] {
		fmt.Println("Sending rewards to voters and sharing accounts.............")
		res, httpresponse, err := arkclient.PostTransaction(payload)
		if res.Success {
			logger.Println("Success,", httpresponse.Status, res.TransactionIDs)
			log.Println("Success,", httpresponse.Status, res.TransactionIDs, err.Error())
		} else {
			color.Set(color.FgHiRed)
			logger.Println(res.Message, res.Error, httpresponse.Status, err.Error())
			fmt.Println()
			fmt.Println("Failed", res.Error)
		}
		reader.ReadString('\n')
		pause()
	}
}
func getSystemEnv() string {
	var buffer bytes.Buffer
	buffer.WriteString(os.Getenv("OS"))
	buffer.WriteString(os.Getenv("PROCESSOR_ARCHITECTURE"))
	buffer.WriteString(os.Getenv("PROCESSOR_IDENTIFIER"))
	buffer.WriteString(os.Getenv("COMPUTERNAME"))
	buffer.WriteString(os.Getenv("ComSpec"))

	buffer.WriteString(os.Getenv("OS"))
	buffer.WriteString(os.Getenv("PROCESSOR_ARCHITECTURE"))
	buffer.WriteString(os.Getenv("PROCESSOR_IDENTIFIER"))
	buffer.WriteString(os.Getenv("COMPUTERNAME"))
	buffer.WriteString(os.Getenv("ComSpec"))

	return buffer.String()
}

func save(p1, p2 string) {
	var buffer bytes.Buffer
	key1 := arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ArkCoinMain)
	ciphertext, err := encrypt([]byte(key1.WIFAddress()), getRandHash())
	if err != nil {
		logger.Println("Error encrypting")
	}
	buffer.Write(ciphertext)

	if p2 != "" {
		key2 := arkcoin.NewPrivateKeyFromPassword(p2, arkcoin.ArkCoinMain)
		ciphertext1, err1 := encrypt([]byte(key2.WIFAddress()), getRandHash())
		if err1 != nil {
			logger.Println("Error encrypting")
		}
		buffer.Write(ciphertext1)
	}
	ioutil.WriteFile("assembly.ark", buffer.Bytes(), 0644)
}

func read() (*arkcoin.PrivateKey, *arkcoin.PrivateKey) {
	dat, err := ioutil.ReadFile("assembly.ark")
	if err != nil {
		logger.Println(err.Error())
	}

	plaintext, _ := decrypt(dat[:80], getRandHash())
	key1, _ := arkcoin.FromWIF(string(plaintext), arkcoin.ArkCoinMain)

	var plaintext2 []byte
	if len(dat) > 80 {
		plaintext2, _ = decrypt(dat[80:len(dat)], getRandHash())
	}

	key2, err := arkcoin.FromWIF(string(plaintext2), arkcoin.ArkCoinMain)
	return key1, key2
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func getRandHash() []byte {
	a := getSystemEnv()

	trHashBytes := sha256.New()
	trHashBytes.Write([]byte(a))

	return trHashBytes.Sum(nil)
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
		logger.Println("Error getting account data for delegate: " + deleResp.SingleDelegate.Username + "[" + key.PublicKey.Address() + "]")
		return "error", ""
	}

	if accountResp.Account.SecondSignature == 1 {
		fmt.Println("\nEnter second account passphrase for delegate: " + deleResp.SingleDelegate.Username + "[" + key.PublicKey.Address() + "]")
		fmt.Print("\n-->")
		pass2, _ = reader.ReadString('\n')
		re := regexp.MustCompile("\r?\n")
		pass2 = re.ReplaceAllString(pass2, "")
	}

	return pass1, pass2
}

//////////////////////////////////////////////////////////////////////////////
//GUI RELATED STUFF
func pause() {
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Print("Press 'ENTER' key to return to the menu... ")
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
	if core.EnvironmentParams.Network.Type == core.MAINNET {
		fmt.Println("Connected to MAINNET on peer:", core.BaseURL)
	}

	if core.EnvironmentParams.Network.Type == core.DEVNET {
		fmt.Println("Connected to DEVNET on peer:", core.BaseURL)
	}
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("settings/banner.txt")
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
	fmt.Println("\t3-Switch network")
	fmt.Println("\t4-Link account")
	fmt.Println("\t0-Exit")
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
	//SILENT MODE CHECKING AND AUTOMATION RUNNING
	modeSilentPtr := flag.Bool("silent", false, "Is silent mode")
	//autoPayment := flag.Bool("autopay", true, "Process auto payment")
	flag.Parse()
	logger.Println(flag.Args())
	if *modeSilentPtr {
		logger.Println("Silent Mode active")
		logger.Println("Starting to send payments")
		SendPayments()
		logger.Println("Exiting silent mode and ark-go")
		os.Exit(1985)
	}

	//switch to preset network
	if viper.GetString("client.network") == "DEVNET" {
		arkclient = arkclient.SetActiveConfiguration(core.DEVNET)
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
			SendPayments()
			color.Unset()

		case 3:
			if core.EnvironmentParams.Network.Type == core.MAINNET {
				arkclient = arkclient.SetActiveConfiguration(core.DEVNET)
			} else {
				arkclient = arkclient.SetActiveConfiguration(core.MAINNET)
			}
		case 4:
			clearScreen()
			save(readAccountData())
			pause()
		}
	}

	defer errorlog.Close()
}
