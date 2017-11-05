package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/kristjank/ark-go/core"
	log "github.com/sirupsen/logrus"
)

//////////////////////////////////////////////////////////////////////////////
//GUI RELATED STUFF
func pause() {
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Print("Press 'ENTER' key to continue... ")
	//bufio.NewReader(os.Stdin).ReadBytes('\n')
	ConsoleReader.ReadString('\n')
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
		log.Info("Connected to ARK MAINNET on peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)
	}

	if core.EnvironmentParams.Network.Type == core.DEVNET {
		fmt.Println("Connected to ARK DEVNET on peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)
		log.Info("Connected to ARK DEVNET on peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)
	}

	if core.EnvironmentParams.Network.Type == core.KAPU {
		fmt.Println("Connected to KAPU MAINNET on peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)
		log.Info("Connected to KAPU MAINNET on peer:", core.BaseURL, "| ARKGoTester version", ArkGoTesterVersion)
	}
}

func printBanner() {
	color.Set(color.FgHiGreen)
	dat, _ := ioutil.ReadFile("cfg/banner.txt")
	fmt.Print(string(dat))
}

func printMenu() {
	log.Info("--------- MAIN MENU ----------------")
	clearScreen()
	printBanner()
	printNetworkInfo()
	color.Set(color.FgHiYellow)
	fmt.Println("")
	fmt.Println("\t1-Run Tests")
	fmt.Println("\t8-Check confirmations")
	fmt.Println("\t9-List DB tests")
	fmt.Println("\t0-Exit")
	fmt.Println("")
	fmt.Print("\tSelect option [0-9]:")
	color.Unset()
}
