package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"wallet/internal/currencies/eth"
	"wallet/internal/utils"
)

func createWallet(password string) *utils.Wallet {
	wallet, _, _ := utils.CreateWallet(password)
	fmt.Println("Wallet created")
	return wallet
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var wallet *utils.Wallet = nil
	var ethAccount *eth.ETHAccount = nil

	for {
		fmt.Println("Enter a command or (type 'exit' to quit):")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "create-wallet":
			fmt.Println("Enter password:")
			scanner.Scan()
			password := strings.TrimSpace(scanner.Text())
			wallet = createWallet(password)
			fmt.Println("Wallet created successfully")
		case "create-eth-account":
			if wallet == nil {
				fmt.Println("Wallet not found")
				break
			}
			ethAccount, _ = wallet.CreateETHAccount()
			fmt.Println("Account created")
			fmt.Println("ETH Address: ", ethAccount.GetAddress())

		case "exit":
			fmt.Println("Exiting...")
			return
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			break
		}
	}
}
