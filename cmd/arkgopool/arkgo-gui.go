package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
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

	"github.com/kristjank/ark-go/arkcoin"
	"github.com/kristjank/ark-go/core"

	"github.com/fatih/color"
	"github.com/spf13/viper"

	"github.com/asdine/storm"
)

var arkclient = core.NewArkClient(nil)
var reader = bufio.NewReader(os.Stdin)

var errorlog *os.File
var logger *log.Logger

var arkpooldb *storm.DB

//PaymentLogRecord structure
type PaymentLogRecord struct {
	Pk              int    `storm:"id,increment"` // primary key with auto increment
	Address         string `storm:"index"`
	VoteWeight      float64
	VoteWeightShare float64
	EarnedAmount100 float64
	EarnedAmountXX  float64
	VoteDuration    int
	Transaction     core.Transaction
}

func save2db(ve core.DelegateDataProfit, tx *core.Transaction) {
	dbData := PaymentLogRecord{}

	dbData.Address = ve.Address
	dbData.VoteWeight = ve.VoteWeight
	dbData.VoteWeightShare = ve.VoteWeightShare
	dbData.EarnedAmount100 = ve.EarnedAmount100
	dbData.EarnedAmountXX = ve.EarnedAmountXX
	dbData.VoteDuration = ve.VoteDuration
	dbData.Transaction = *tx

	err := arkpooldb.Save(&dbData)
	if err != nil {
		log.Println(err.Error())
	}
	pause()
}

func listPaymentsFromDB() {
	var results []PaymentLogRecord
	err := arkpooldb.All(&results)

	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, element := range results {
		log.Println(element.Transaction.RecipientID)
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
func init() {
	var err error
	errorlog, err = os.OpenFile("arkgo-gui.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
		os.Exit(1)
	}

	logger = log.New(errorlog, "ark-go: ", log.Lshortfile|log.LstdFlags)
}

func initializeBoltClient() {
	var err error
	arkpooldb, err = storm.Open(viper.GetString("client.dbfilename"))

	if err != nil {
		panic(err.Error())
	}

	log.Println("DB Opened at:", arkpooldb.Path)
	//defer arkpooldb.Close()
}

//DisplayCalculatedVoteRatio based on parameters in config.toml
func DisplayCalculatedVoteRatio() {
	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	var key1 *arkcoin.PrivateKey
	var p1 string
	isLinked := false
	if _, err := os.Stat("assembly.ark"); err == nil {
		logger.Println("Linked accound data found. Using saved account information.")
		p1, _ = read()
		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
		pubKey = hex.EncodeToString(key1.PublicKey.Serialize())
		isLinked = true
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}
	deleResp, _, _ := arkclient.GetDelegate(params)
	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"))
	shareRatioStr := strconv.FormatFloat(viper.GetFloat64("voters.shareratio")*100, 'f', -1, 64) + "%"

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("Displaying voter information for delegate:")
	color.Set(color.FgHiYellow)
	fmt.Print("\tusername:", deleResp.SingleDelegate.Username)
	fmt.Println("\taddress:", deleResp.SingleDelegate.Address)
	fmt.Print("\tfidelity:")
	color.HiRed("%t", viper.GetBool("voters.fidelity"))
	color.Set(color.FgHiYellow)
	fmt.Print("\tfee deduction:")
	color.HiRed("%t", viper.GetBool("voters.deductTxFees"))
	color.Set(color.FgHiYellow)
	fmt.Print("\tlinked:")
	color.HiRed("%t\n", isLinked)
	color.Set(color.FgHiGreen)

	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("%4s|%-34s|%18s|%8s|%17s|%6s|%15s|", "Ix", "Voter address", "Balance", "Weight", "Reward-"+shareRatioStr, "Hours", "FidelityAmount"))
	color.Set(color.FgCyan)
	for ix, element := range votersEarnings {

		fidelAmount := calcFidelity(element)

		s := fmt.Sprintf("%3d.|%s|%18.8f|%8.4f|%15.8f A|%6d|%15.8f|", ix+1, element.Address, element.VoteWeight, element.VoteWeightShare, element.EarnedAmountXX, element.VoteDuration, fidelAmount)

		fmt.Println(s)
		logger.Println(s)

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}

	//Cost calculation
	costAmount := sumEarned * viper.GetFloat64("costs.shareratio")
	reserveAmount := sumEarned * viper.GetFloat64("reserve.shareratio")
	personalAmount := sumEarned * viper.GetFloat64("personal.shareratio")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("Available amount:", sumEarned)
	fmt.Println("Amount to voters:", sumShareEarned, viper.GetFloat64("voters.shareratio"))
	fmt.Println("Amount to costs:", costAmount, viper.GetFloat64("costs.shareratio"))
	fmt.Println("Amount to reserve:", reserveAmount, viper.GetFloat64("reserve.shareratio"))
	fmt.Println("Amount to personal:", personalAmount, viper.GetFloat64("personal.shareratio"))

	fmt.Println("Ratio calc check:", sumRatio, "(should be = 1)")
	fmt.Println("Ratio share check:", float64(sumShareEarned)/float64(sumEarned), "should be=", viper.GetFloat64("voters.shareratio"))

	pause()
}

func floatEquals(a, b float64) bool {
	EPSILON := 0.000000000000001
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	}
	return false
}

func checkConfigSharingRatio() bool {
	a1 := viper.GetFloat64("voters.shareratio")
	a2 := viper.GetFloat64("costs.shareratio")
	a3 := viper.GetFloat64("reserve.shareratio")
	a4 := viper.GetFloat64("personal.shareratio")

	if !floatEquals(a1+a2+a3+a4, 1.0) {
		logger.Println("Wrong config. Check share ration percentages!")
		return false
	}
	return true
}

//SendPayments based on parameters in config.toml
func SendPayments(silent bool) {
	if !checkConfigSharingRatio() {
		clearScreen()
		color.Set(color.FgHiRed)
		fmt.Println("--------------------------------------------------------------------------------------------------------------")
		fmt.Println("")
		fmt.Println("Unable to calculate. Check share ratio configuration.")
		pause()
		logger.Println("Unable to calculcate. Check share ratio configuration.")
		return
	}

	isLinked := false
	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	var p1, p2 string
	var key1 *arkcoin.PrivateKey
	if _, err := os.Stat("assembly.ark"); err == nil {
		logger.Println("Linked accound data found. Using saved account information.")

		p1, p2 = read()

		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
		pubKey = hex.EncodeToString(key1.PublicKey.Serialize())
		isLinked = true
	} else {
		p1, p2 = readAccountData()
		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}
	var payload core.TransactionPayload

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"))

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := 0

	clearScreen()

	for _, element := range votersEarnings {
		//Logging history to DB

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		fAmount2Send := calcFidelity(element)

		//transaction parameters
		txAmount2Send := int64(fAmount2Send * core.SATOSHI)

		//decuting fees if setup
		if viper.GetBool("voters.deductTxFees") {
			txAmount2Send -= core.EnvironmentParams.Fees.Send
			logger.Println("Voters Fee deduction enabled")
		}

		//only payout for earning higher then minamount. - the earned amount remains in the loop for next payment
		//to disable set it to 0.0
		if element.EarnedAmountXX >= viper.GetFloat64("voters.minamount") && txAmount2Send > 0 {
			tx := core.CreateTransaction(element.Address, txAmount2Send, viper.GetString("voters.txdescription"), p1, p2)
			payload.Transactions = append(payload.Transactions, tx)

			//Logging history to DB
			save2db(element, tx)
		}
	}

	//if decuting fees from voters is false - we take them into account here....
	//must be at this spot - as it counts the number of voters to get the rewards - befor other
	//transactions are added...
	if !viper.GetBool("voters.deductTxFees") {
		feeAmount = int(len(payload.Transactions)) * int(core.EnvironmentParams.Fees.Send)
	}

	//Cost & reserve fund calculation
	costAmount := sumEarned * viper.GetFloat64("costs.shareratio")
	reserveAmount := sumEarned * viper.GetFloat64("reserve.shareratio")
	personalAmount := sumEarned * viper.GetFloat64("personal.shareratio")

	//summary and conversion checks
	if (costAmount + reserveAmount + personalAmount + sumShareEarned) != sumEarned {
		color.Set(color.FgHiRed)
		diff := sumEarned - (costAmount + reserveAmount + personalAmount + sumShareEarned)
		if diff > 0.00000001 {
			log.Println("Calculation of ratios NOT OK - overall summary failing for diff=", diff)
			logger.Println("Calculation of ratios NOT OK - overall summary failing diff=", diff)
		}
	}

	//cost amount calculation
	costAmount2Send := int64(costAmount*core.SATOSHI) - core.EnvironmentParams.Fees.Send
	if costAmount2Send > 0 {
		costAddress := viper.GetString("costs.address")
		if core.EnvironmentParams.Network.Type == core.DEVNET {
			costAddress = viper.GetString("costs.Daddress")
		}

		txCosts := core.CreateTransaction(costAddress, costAmount2Send, viper.GetString("costs.txdescription"), p1, p2)
		payload.Transactions = append(payload.Transactions, txCosts)
	}

	//Reserve amount
	reserveAmount2Send := int64(reserveAmount*core.SATOSHI) - core.EnvironmentParams.Fees.Send
	if reserveAmount2Send > 0 {
		reserveAddress := viper.GetString("reserve.address")
		if core.EnvironmentParams.Network.Type == core.DEVNET {
			reserveAddress = viper.GetString("reserve.Daddress")
		}
		txReserve := core.CreateTransaction(reserveAddress, reserveAmount2Send, viper.GetString("reserve.txdescription"), p1, p2)
		payload.Transactions = append(payload.Transactions, txReserve)
	}

	//Personal
	personalAmount2Send := int64(personalAmount*core.SATOSHI) - core.EnvironmentParams.Fees.Send
	if personalAmount2Send > 0 {
		personalAddress := viper.GetString("personal.address")
		if core.EnvironmentParams.Network.Type == core.DEVNET {
			personalAddress = viper.GetString("personal.Daddress")
		}
		txpersonal := core.CreateTransaction(personalAddress, personalAmount2Send, viper.GetString("personal.txdescription"), p1, p2)
		payload.Transactions = append(payload.Transactions, txpersonal)
	}

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("Transactions to be sent from:")
	color.Set(color.FgHiYellow)
	fmt.Println("\tDelegate address:", key1.PublicKey.Address())
	color.Set(color.FgHiYellow)
	fmt.Print("\tFidelity:")
	color.HiRed("%t", viper.GetBool("voters.fidelity"))
	color.Set(color.FgHiYellow)
	fmt.Print("\tFee deduction:")
	color.HiRed("%t", viper.GetBool("voters.deductTxFees"))
	color.Set(color.FgHiYellow)
	fmt.Println("\tFee Amount:", feeAmount)
	color.Set(color.FgHiYellow)
	fmt.Print("\tLinked:")
	color.HiRed("%t\n", isLinked)
	color.Set(color.FgHiGreen)

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	color.Set(color.FgHiCyan)
	for ix, el := range payload.Transactions {
		s := fmt.Sprintf("%3d.|%s|%15d| %-40s|", ix+1, el.RecipientID, el.Amount, el.VendorField)
		fmt.Println(s)
		logger.Println(s)
	}

	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")

	var c byte
	if !silent {
		fmt.Print("Send transactions and complete reward payments [Y/N]: ")
		c, _ = reader.ReadByte()
	} else {
		fmt.Print("Sending automated transactions")
		logger.Println("Sending automated transactions")
		c = []byte("Y")[0]
	}

	if c == []byte("Y")[0] || c == []byte("y")[0] {
		fmt.Println("Sending rewards to voters and sharing accounts.............")

		res, httpresponse, err := arkclient.PostTransaction(payload)
		if res.Success {
			color.Set(color.FgHiGreen)
			logger.Println("Transactions sent with Success,", httpresponse.Status, res.TransactionIDs)
			log.Println("Transactions sent with Success,", httpresponse.Status)
			log.Println("Audit log of sent transactions is in file paymentLog.csv!")
			log2csv(payload, res.TransactionIDs, votersEarnings)
		} else {
			color.Set(color.FgHiRed)
			logger.Println(res.Message, res.Error, httpresponse.Status, err.Error())
			fmt.Println()
			fmt.Println("Failed", res.Error)
		}
		if !silent {
			reader.ReadString('\n')
			pause()
		}
	}
}

//SendBonus Send fixed amount to all voters
func SendBonus() {

	//Bonus to send
	fmt.Println("\nEnter amount to send to each voter")
	fmt.Print("-->")
	bonusInput, _ := reader.ReadString('\n')
	re := regexp.MustCompile("\r?\n")
	bonusInput = re.ReplaceAllString(bonusInput, "")

	bonus, err := strconv.ParseFloat(bonusInput, 64)
	if err != nil {
		panic(err)
	}

	//Bonus description
	fmt.Println("\nEnter transaction description")
	fmt.Print("-->")
	txDescription, _ := reader.ReadString('\n')
	res := regexp.MustCompile("\r?\n")
	txDescription = res.ReplaceAllString(txDescription, "")

	isLinked := false
	pubKey := viper.GetString("delegate.pubkey")
	if core.EnvironmentParams.Network.Type == core.DEVNET {
		pubKey = viper.GetString("delegate.Dpubkey")
	}

	var p1, p2 string
	var key1 *arkcoin.PrivateKey
	if _, err := os.Stat("assembly.ark"); err == nil {
		logger.Println("Linked accound data found. Using saved account information.")

		p1, p2 = read()

		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
		pubKey = hex.EncodeToString(key1.PublicKey.Serialize())
		isLinked = true
	} else {
		p1, p2 = readAccountData()
		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
	}

	params := core.DelegateQueryParams{PublicKey: pubKey}
	var payload core.TransactionPayload

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"))

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0

	clearScreen()

	for _, element := range votersEarnings {
		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		//Bonus amount
		txAmount2Send := int64(bonus * core.SATOSHI)

		//decuting fees if setup
		if viper.GetBool("voters.deductTxFees") {
			txAmount2Send -= core.EnvironmentParams.Fees.Send
			logger.Println("Voters Fee deduction enabled")
		}

		//Bonus - only pay account with more than 0 balance
		if element.EarnedAmount100 > 0 && txAmount2Send > 0 {
			tx := core.CreateTransaction(element.Address, txAmount2Send, txDescription, p1, p2)
			payload.Transactions = append(payload.Transactions, tx)
		}
	}

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("Transactions to be sent from:")
	color.Set(color.FgHiYellow)
	fmt.Println("\tDelegate address:", key1.PublicKey.Address(), "linked:", isLinked)
	color.Set(color.FgHiYellow)
	fmt.Print("\tFidelity:")
	color.HiRed("%t", viper.GetBool("voters.fidelity"))
	color.Set(color.FgHiYellow)
	fmt.Print("\tFee deduction:")
	color.HiRed("%t", viper.GetBool("voters.deductTxFees"))
	color.Set(color.FgHiYellow)
	fmt.Print("\tLinked:")
	color.HiRed("%t\n", isLinked)
	color.Set(color.FgHiGreen)

	color.Set(color.FgHiGreen)
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	color.Set(color.FgHiCyan)
	for ix, el := range payload.Transactions {
		s := fmt.Sprintf("%3d.|%s|%15d| %-40s|", ix+1, el.RecipientID, el.Amount, el.VendorField)
		fmt.Println(s)
		logger.Println(s)
	}

	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("--------------------------------------------------------------------------------------------------------------")

	var c byte
	fmt.Print("Send transactions and complete bonus payments [Y/N]: ")
	c, _ = reader.ReadByte()

	if c == []byte("Y")[0] || c == []byte("y")[0] {
		fmt.Println("Sending bonus to voters accounts.............")

		res, httpresponse, err := arkclient.PostTransaction(payload)
		if res.Success {
			color.Set(color.FgHiGreen)
			logger.Println("Transactions sent with Success,", httpresponse.Status, res.TransactionIDs)
			log.Println("Transactions sent with Success,", httpresponse.Status)
			log.Println("Audit log of sent transactions is in file paymentLog.csv!")
			log2csv(payload, res.TransactionIDs, votersEarnings)
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

func calcFidelity(element core.DelegateDataProfit) float64 {
	fAmount2Send := element.EarnedAmountXX
	//FIDELITY
	if viper.GetBool("voters.fidelity") {
		if element.VoteDuration < viper.GetInt("voters.fidelityLimit") {
			fAmount2Send *= float64(element.VoteDuration) / float64(viper.GetInt("voters.fidelityLimit"))
			logger.Println("Fidelity enabled for user", element.Address, "ratio: ", float64(element.VoteDuration)/float64(viper.GetInt("voters.fidelityLimit")), "earned: ", element.EarnedAmountXX, "reduced amount: ", fAmount2Send)
		}
	}

	return fAmount2Send
}

func log2csv(payload core.TransactionPayload, txids []string, voterCalcs []core.DelegateDataProfit) {
	records := [][]string{
		{"ADDRESS", "SENT AMOUNT", "WALLET BALANCE", "Fidelity(h)", "TimeStamp", "TxId"},
	}

	for ix, el := range payload.Transactions {
		//		sAmount := fmt.Sprintf("%15.8f", float64(el.Amount)/float64(core.SATOSHI))
		timeTx := core.GetTransactionTime(el.Timestamp)
		localTime := timeTx.Local()

		wBalance := "N/A"
		wDuration := "N/A"
		if ix < len(voterCalcs) {
			wBalance = strconv.FormatFloat(voterCalcs[ix].VoteWeight, 'f', -1, 64)
			wDuration = strconv.FormatInt(int64(voterCalcs[ix].VoteDuration), 10)
		}

		line := []string{el.RecipientID, strconv.FormatFloat(float64(el.Amount)/float64(core.SATOSHI), 'f', -1, 64), wBalance, wDuration, localTime.Format("2006-01-02 15:04:05"), txids[ix]}
		records = append(records, line)

	}
	file, _ := os.Create("paymentLog.csv")
	w := csv.NewWriter(file)
	defer w.Flush()
	w.WriteAll(records)
	file.Close()
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
	ciphertext, _ := encrypt([]byte(p1), getRandHash())
	ioutil.WriteFile("assembly.ark", ciphertext, 0644)

	if p2 != "" {
		ciphertext, err := encrypt([]byte(p2), getRandHash())
		if err != nil {
			logger.Println("Error encrypting")
		}
		ioutil.WriteFile("assembly1.ark", ciphertext, 0644)
	} else {
		os.Remove("assembly1.ark")
	}
}

/*func read() (*arkcoin.PrivateKey, *arkcoin.PrivateKey) {
	dat, err := ioutil.ReadFile("assembly.ark")
	if err != nil {
		logger.Println(err.Error())
	}
	plaintext, _ := decrypt(dat, getRandHash())
	key1 := arkcoin.NewPrivateKeyFromPassword(string(plaintext), arkcoin.ActiveCoinConfig)

	var key2 *arkcoin.PrivateKey
	if _, err := os.Stat("assembly1.ark"); err == nil {
		dat, _ = ioutil.ReadFile("assembly1.ark")
		plaintext, _ = decrypt(dat, getRandHash())
		key2 = arkcoin.NewPrivateKeyFromPassword(string(plaintext), arkcoin.ActiveCoinConfig)
	}

	return key1, key2
}*/

func read() (string, string) {
	dat, err := ioutil.ReadFile("assembly.ark")
	if err != nil {
		logger.Println(err.Error())
	}
	p1, _ := decrypt(dat, getRandHash())

	var p2 []byte

	if _, err := os.Stat("assembly1.ark"); err == nil {
		dat, _ = ioutil.ReadFile("assembly1.ark")
		p2, _ = decrypt(dat, getRandHash())
	}

	return string(p1), string(p2)
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
		fmt.Print("-->")
		pass2, _ = reader.ReadString('\n')
		re := regexp.MustCompile("\r?\n")
		pass2 = re.ReplaceAllString(pass2, "")
	}

	return pass1, pass2
}

func loadConfig() {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.AddConfigPath("settings") // path to look for the config file in
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file

	if err != nil {
		logger.Println("No productive config found - loading sample")
		// try to load sample config
		viper.SetConfigName("sample.config")
		viper.AddConfigPath("settings")
		err := viper.ReadInConfig()

		if err != nil { // Handle errors reading the config file
			logger.Println("No configuration file loaded - using defaults")
		}
	}

	viper.SetDefault("delegate.address", "")
	viper.SetDefault("delegate.pubkey", "")
	viper.SetDefault("delegate.Daddress", "")
	viper.SetDefault("delegate.Dpubkey", "")

	viper.SetDefault("voters.shareRatio", 0.0)
	viper.SetDefault("voters.txdescription", "share tx by ark-go")
	viper.SetDefault("voters.fidelity", true)
	viper.SetDefault("voters.fidelityLimit", 24)
	viper.SetDefault("voters.minamount", 0.0)
	viper.SetDefault("voters.deductTxFees", true)

	viper.SetDefault("costs.address", "")
	viper.SetDefault("costs.shareRatio", 0.0)
	viper.SetDefault("costs.txdescription", "cost tx by ark-go")
	viper.SetDefault("costs.Daddress", "")

	viper.SetDefault("reserve.address", "")
	viper.SetDefault("reserve.shareRatio", 0.0)
	viper.SetDefault("reserve.txdescription", "reserve tx by ark-go")
	viper.SetDefault("reserve.Daddress", "")

	viper.SetDefault("personal.address", "")
	viper.SetDefault("personal.shareRatio", 0.0)
	viper.SetDefault("personal.txdescription", "personal tx by ark-go")
	viper.SetDefault("personal.Daddress", "")

	viper.SetDefault("client.network", "DEVNET")
	viper.SetDefault("client.dbFilename", "payment.db")
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
	fmt.Println("\t2-Send reward payments")
	fmt.Println("\t3-Switch network")
	fmt.Println("\t4-Link account")
	fmt.Println("\t5-Send bonus payments")
	fmt.Println("\t6-List history payments")
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

	// Load configration and defaults
	loadConfig()

	initializeBoltClient()

	//switch to preset network
	if viper.GetString("client.network") == "DEVNET" {
		arkclient = arkclient.SetActiveConfiguration(core.DEVNET)
	}

	//SILENT MODE CHECKING AND AUTOMATION RUNNING
	modeSilentPtr := flag.Bool("silent", false, "Is silent mode")
	//autoPayment := flag.Bool("autopay", true, "Process auto payment")
	flag.Parse()
	logger.Println(flag.Args())
	if *modeSilentPtr {
		logger.Println("Silent Mode active")
		logger.Println("Starting to send payments")
		SendPayments(true)
		logger.Println("Exiting silent mode and ark-go")
		color.Unset()
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
			color.Set(color.FgHiGreen)
			logger.Println("Account succesfully linked")
			fmt.Println("Account succesfully linked")
			pause()
			color.Unset()
		case 5:
			clearScreen()
			color.Set(color.FgHiGreen)
			SendBonus()
			color.Unset()
		case 6:
			clearScreen()
			color.Set(color.FgHiGreen)
			listPaymentsFromDB()
			pause()
			color.Unset()
		}
	}
	color.Unset()
	defer errorlog.Close()
}
