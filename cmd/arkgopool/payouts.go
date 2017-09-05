package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kristjank/ark-go/arkcoin"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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
		log.Info("Linked accound data found. Using saved account information.")
		p1, _ = read()
		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
		pubKey = hex.EncodeToString(key1.PublicKey.Serialize())
		isLinked = true
	}

	//TODO REMOVE BEFORE build
	//pubKey = "02c7455bebeadde04728441e0f57f82f972155c088252bf7c1365eb0dc84fbf5de"

	params := core.DelegateQueryParams{PublicKey: pubKey}
	deleResp, _, _ := arkclient.GetDelegate(params)
	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), viper.GetString("voters.blocklist"))
	shareRatioStr := strconv.FormatFloat(viper.GetFloat64("voters.shareratio")*100, 'f', -1, 64) + "%"

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := 0.0

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
		log.Info(s)

		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare
	}

	//Cost calculation
	costAmount := sumEarned * viper.GetFloat64("costs.shareratio")
	reserveAmount := sumEarned * viper.GetFloat64("reserve.shareratio")
	personalAmount := sumEarned * viper.GetFloat64("personal.shareratio")

	//if decuting fees from voters is false - we take them into account here....
	//must be at this spot - as it counts the number of voters to get the rewards - befor other
	//transactions are added...
	if !viper.GetBool("voters.deductTxFees") {
		feeAmount = float64(int(len(votersEarnings))*int(core.EnvironmentParams.Fees.Send)) / float64(core.SATOSHI)
		log.Info("Calculated fee amount: ", feeAmount)

		//deducting feeAmount from reserve address
		if feeAmount > reserveAmount {
			log.Fatal("Not enough reserve money to pay the fees from reserve fund. Payment script stopped !!!")
		}
		reserveAmount -= float64(feeAmount)
	}
	fmt.Println("--------------------------------------------------------------------------------------------------------------")
	fmt.Println("")
	fmt.Println("Available amount:", sumEarned)
	fmt.Println("Amount to voters:", sumShareEarned, viper.GetFloat64("voters.shareratio"))
	fmt.Println("Amount to costs:", costAmount, viper.GetFloat64("costs.shareratio"))
	fmt.Println("Amount to reserve:", reserveAmount, viper.GetFloat64("reserve.shareratio"))
	fmt.Println("Amount for fees:", feeAmount, " if bigger > 0, it is deducted from reserve amount")
	fmt.Println("Fee calc check (reserve2send+fee):", feeAmount+reserveAmount, " should be=", sumEarned*viper.GetFloat64("reserve.shareratio"))
	fmt.Println("Amount to personal:", personalAmount, " share ratio: ", viper.GetFloat64("personal.shareratio"))
	fmt.Println("Voters Ratio calc check:", sumRatio, " (should be = 1)")
	fmt.Println("Voters Ratio share check:", float64(sumShareEarned)/float64(sumEarned), " should be=", viper.GetFloat64("voters.shareratio"))

	pause()
}

//SendPayments based on parameters in config.toml
func SendPayments(silent bool) {
	payrec := createPaymentRecord()
	arkpooldb.Save(&payrec)
	log.Info("Starting payments calculation. Active peer for voter information: ", arkclient.GetActivePeer())

	if !checkConfigSharingRatio() {
		clearScreen()
		color.Set(color.FgHiRed)
		if !silent {
			fmt.Println("--------------------------------------------------------------------------------------------------------------")
			fmt.Println("")
			fmt.Println("Unable to calculate. Check share ratio configuration.")
			pause()
		}
		log.Fatal("Unable to calculcate. Check share ratio configuration in your config.toml.")
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
		log.Info("Linked accound data found. Using saved account information.")

		p1, p2 = read()

		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
		pubKey = hex.EncodeToString(key1.PublicKey.Serialize())

		isLinked = true
	} else {
		p1, p2 = readAccountData()
		key1 = arkcoin.NewPrivateKeyFromPassword(p1, arkcoin.ActiveCoinConfig)
	}

	//TODO JARUNIK TEST
	//pubKey = "02c7455bebeadde04728441e0f57f82f972155c088252bf7c1365eb0dc84fbf5de"
	//pubKey = "027acdf24b004a7b1e6be2adf746e3233ce034dbb7e83d4a900f367efc4abd0f21"
	params := core.DelegateQueryParams{PublicKey: pubKey}
	var payload core.TransactionPayload

	deleResp, _, _ := arkclient.GetDelegateVoters(params)

	// check minVoteDuration
	var blocklist = checkMinimumVoteTime(deleResp, viper.GetString("voters.blocklist"))

	votersEarnings := arkclient.CalculateVotersProfit(params, viper.GetFloat64("voters.shareratio"), blocklist)
	payrec.VoteWeight, _, _ = arkclient.GetDelegateVoteWeight(params)

	sumEarned := 0.0
	sumRatio := 0.0
	sumShareEarned := 0.0
	feeAmount := 0.0
	minAmountSetting := int64(viper.GetFloat64("voters.minamount") * core.SATOSHI)

	clearScreen()

	//calculating voter earnings
	for _, element := range votersEarnings {
		sumEarned += element.EarnedAmount100
		sumShareEarned += element.EarnedAmountXX
		sumRatio += element.VoteWeightShare

		fAmount2Send := calcFidelity(element)

		//transaction parameters
		txAmount2Send := int64(fAmount2Send * core.SATOSHI)

		//decuting fees if setup
		if viper.GetBool("voters.deductTxFees") {
			txAmount2Send -= core.EnvironmentParams.Fees.Send
			log.Info("Voters Fee deduction enabled")
		}

		//checking MinAmount && MaxAmount properties
		if txAmount2Send > minAmountSetting && txAmount2Send > 0 {
			tx := core.CreateTransaction(element.Address, txAmount2Send, viper.GetString("voters.txdescription"), p1, p2)
			payload.Transactions = append(payload.Transactions, tx)
			//Logging history to DB
			save2db(element, tx, payrec.Pk)
		} else {
			log.Info("Skipping voter address ", element.Address, " Earned amount: ", txAmount2Send, " below minimium: ", minAmountSetting)
		}
	}

	//Cost & reserve fund calculation
	costAmount := sumEarned * viper.GetFloat64("costs.shareratio")
	reserveAmount := sumEarned * viper.GetFloat64("reserve.shareratio")
	personalAmount := sumEarned * viper.GetFloat64("personal.shareratio")

	//if decuting fees from voters is false - we take them into account here....
	//must be at this spot - as it counts the number of voters to get the rewards - befor other
	//transactions are added, and only voters with enough big share to payout
	if !viper.GetBool("voters.deductTxFees") {
		feeAmount = float64(int(len(payload.Transactions))*int(core.EnvironmentParams.Fees.Send)) / float64(core.SATOSHI)
		log.Info("Calculated fee amount: ", feeAmount)
		payrec.FeeAmount = feeAmount

		//deducting feeAmount from reserve address
		if feeAmount > reserveAmount {
			log.Fatal("Not enough reserve money to pay the fees from reserve fund. Payment script stopped !!!")
		}
		reserveAmount -= float64(feeAmount)
	}
	log.Info("*******************************************************************************************************************")
	log.Info("                                CALCULATION INFORMATION FOR THIS RUN ")
	log.Info("Available amount:", sumEarned)
	log.Info("Number of voters (excluding blacklisted):", len(votersEarnings))
	log.Info("Amount to voters:", sumShareEarned, " Share ratio: ", viper.GetFloat64("voters.shareratio"))
	log.Info("Amount to costs:", costAmount, " Share ratio: ", viper.GetFloat64("costs.shareratio"))
	log.Info("Amount to reserve:", reserveAmount, " Share ratio: ", viper.GetFloat64("reserve.shareratio"), "deducted by tx fees")
	log.Info("Amount for fees:", feeAmount, " if bigger > 0, it is deducted from reserve amount")
	log.Info("Fee calc check (reserve2send+fee):", feeAmount+reserveAmount, " should be=", sumEarned*viper.GetFloat64("reserve.shareratio"))
	log.Info("Amount to personal:", personalAmount, " Share ratio: ", viper.GetFloat64("personal.shareratio"))
	log.Info("Voters Ratio calc check:", sumRatio, "(should be = 1)")
	log.Info("Voters Ratio share check:", float64(sumShareEarned)/float64(sumEarned), "should be=", viper.GetFloat64("voters.shareratio"))
	log.Info()
	log.Info("*******************************************************************************************************************")

	//summary and conversion checks
	if (costAmount + reserveAmount + personalAmount + sumShareEarned + feeAmount) != sumEarned {
		color.Set(color.FgHiRed)
		diff := sumEarned - (costAmount + reserveAmount + personalAmount + sumShareEarned + feeAmount)
		if diff > 0.00000001 {
			fmt.Println("Calculation of ratios NOT OK - overall summary failing for diff=", diff)
			log.Fatal("Calculation of ratios NOT OK - overall summary failing diff=", diff)
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

	payrec.NrOfTransactions = len(payload.Transactions)
	payrec.FeeAmount = float64(int(len(payload.Transactions))*int(core.EnvironmentParams.Fees.Send)) / float64(core.SATOSHI)

	arkpooldb.Update(&payrec)

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
		log.Info(s)
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
		log.Info("Sending automated transactions")
		c = []byte("Y")[0]
	}

	if c == []byte("Y")[0] || c == []byte("y")[0] {

		fmt.Println("Sending rewards to voters and sharing accounts.............")
		log.Info("Starting automated payment... ")

		splitAndDeliverPayload(payload)

		fmt.Println("Automated Payment complete. Please check the logs folder... ")
		log.Info("Automated Payment complete. Please check the logs folder... ")

		if !silent {
			reader.ReadString('\n')
			pause()
		}

	}
}

func splitAndDeliverPayload(payload core.TransactionPayload) {
	//calculating number of chunks (based on 20tx in one chunk to send to one peer)
	payoutsFolderName := createLogFolder()
	var divided [][]*core.Transaction
	numPeers := len(payload.Transactions) / 20
	if numPeers == 0 {
		numPeers = 1
	}
	chunkSize := (len(payload.Transactions) + numPeers - 1) / numPeers
	if chunkSize == 0 {
		chunkSize = 1
	}

	//sliptting the payload to number of needed peers
	for i := 0; i < len(payload.Transactions); i += chunkSize {
		end := i + chunkSize
		if end > len(payload.Transactions) {
			end = len(payload.Transactions)
		}
		divided = append(divided, payload.Transactions[i:end])
	}
	//end of spliting transactions

	var tmpPayload core.TransactionPayload
	splitcout := 0
	for chunkIx, h := range divided {
		tmpPayload.Transactions = h
		splitcout += len(h)

		deliverPayloadThreaded(tmpPayload, chunkIx, payoutsFolderName)

	}
	if splitcout != len(payload.Transactions) {
		log.Info("TX spliting not OK")
	}
}

func deliverPayloadThreaded(tmpPayload core.TransactionPayload, chunkIx int, logFolder string) {
	numberOfPeers2MultiBroadCastTo := viper.GetInt("client.multibroadcast")
	log.Info("Starting multibroadcast/multithreaded parallel payout to ", numberOfPeers2MultiBroadCastTo, " number of peers")
	peers := arkclient.GetRandomXPeers(numberOfPeers2MultiBroadCastTo)
	for i := 0; i < numberOfPeers2MultiBroadCastTo; i++ {
		wg.Add(1)

		//treaded function
		go func(tmpPayload core.TransactionPayload, peer core.Peer, chunkIx int, logFolder string) {
			defer wg.Done()
			filename := fmt.Sprintf("log/%s/Batch_%2d_Peer%s.csv", logFolder, chunkIx, peer.IP)

			arkTmpClient := core.NewArkClientFromPeer(peer)
			res, _, _ := arkTmpClient.PostTransaction(tmpPayload)
			if res.Success {
				color.Set(color.FgHiGreen)
				log2csv(tmpPayload, res.TransactionIDs, filename, "OK")
			} else {
				color.Set(color.FgHiRed)
				log2csv(tmpPayload, nil, filename, res.Error)
			}
		}(tmpPayload, peers[i], chunkIx, logFolder)
	}
}

func calcFidelity(element core.DelegateDataProfit) float64 {
	fAmount2Send := element.EarnedAmountXX
	//FIDELITY
	if viper.GetBool("voters.fidelity") {
		if element.VoteDuration < viper.GetInt("voters.fidelityLimit") {
			fAmount2Send *= float64(element.VoteDuration) / float64(viper.GetInt("voters.fidelityLimit"))
			log.Info("Fidelity enabled for user ", element.Address, " ratio: ", float64(element.VoteDuration)/float64(viper.GetInt("voters.fidelityLimit")), " earned: ", element.EarnedAmountXX, "reduced amount: ", fAmount2Send)
		}
	}

	return fAmount2Send
}

func checkMinimumVoteTime(voters core.DelegateVoters, blocklist string) string {
	var minVoteDuration = viper.GetInt("voters.minVoteDuration")

	for _, element := range voters.Accounts {
		if minVoteDuration > arkclient.GetVoteDuration(element.Address) {
			if len(blocklist) > 0 {
				if !strings.Contains(strings.ToLower(blocklist), strings.ToLower(element.Address)) {
					blocklist += "," + element.Address
				}
			} else {
				blocklist += element.Address
			}
		}
	}

	return blocklist
}
