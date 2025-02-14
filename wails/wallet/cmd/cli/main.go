package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"wallet/internal/currencies/eth"
	"wallet/internal/utils"
)

func createWallet(password string) (*utils.Wallet, error) {
	wallet, _, err := utils.CreateWallet(password)
	if err != nil {
		return nil, fmt.Errorf("error creating wallet: %v", err)
	}

	return wallet, nil
}

func RestoreWallet(password, mnemonic string) (*utils.Wallet, error) {
	wallet, err := utils.RestoreWallet(password, mnemonic)
	if err != nil {
		return nil, fmt.Errorf("error restoring wallet: %v", err)
	}

	return wallet, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var wallet *utils.Wallet = nil
	var ethAccount *eth.ETHAccount = nil
	provider := "http://127.0.0.1:8545"

	for {
		fmt.Println("Enter a command or (type 'exit' to quit):")
		scanner.Scan()
		command := scanner.Text()

		switch command {
		case "create-wallet":
			fmt.Println("Enter password: ")
			scanner.Scan()
			password := strings.TrimSpace(scanner.Text())
			var err error
			wallet, err = createWallet(password)
			if err != nil {
				fmt.Println("Error creating wallet:", err)
				break
			}
			fmt.Println("Wallet created successfully")

		case "restore-wallet":
			fmt.Println("Enter password: ")
			scanner.Scan()
			password := strings.TrimSpace(scanner.Text())
			fmt.Println("Enter mnemonic: ")
			scanner.Scan()
			mnemonic := strings.TrimSpace(scanner.Text())
			fmt.Println("mnemonic: ", mnemonic)
			var err error
			wallet, err = RestoreWallet(password, mnemonic)
			if err != nil {
				fmt.Println("Error restoring wallet:", err)
				break
			}

			fmt.Println("Wallet restored successfully")

		case "create-eth-account":
			if wallet == nil {
				fmt.Println("Wallet not found")
				break
			}
			ethAccount, _ = wallet.CreateETHAccount()
			fmt.Println("Account created")
			fmt.Println("ETH Address: ", ethAccount.GetAddress())

		case "get-eth-balance":
			if ethAccount == nil {
				fmt.Println("ETH account not found")
				break
			}

			if !eth.NetListening(provider) {
				fmt.Println("Node is not listening")
				break
			}

			fmt.Println("Getting ETH balance from address: ", ethAccount.GetAddress())
			hexBalance, err := eth.GetBalance(provider, ethAccount.GetAddress())
			if err != nil {
				fmt.Println("Failed to get balance: ", err)
				break
			}
			balance, _ := eth.HexToEther(hexBalance)
			fmt.Println("ETH Balance: ", balance)

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
