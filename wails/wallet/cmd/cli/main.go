package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"wallet/internal/currencies/eth"
	"wallet/internal/hdwallet"
)

func createWallet(password string) (*hdwallet.Wallet, error) {
	wallet, _, err := hdwallet.CreateWallet(password)
	if err != nil {
		return nil, fmt.Errorf("error creating wallet: %v", err)
	}

	return wallet, nil
}

func RestoreWallet(password, mnemonic string) (*hdwallet.Wallet, error) {
	wallet, err := hdwallet.RestoreWallet(password, mnemonic)
	if err != nil {
		return nil, fmt.Errorf("error restoring wallet: %v", err)
	}

	return wallet, nil
}

func GetBalance(wallet *hdwallet.Wallet, token, network string) (string, error) {
	balanceFloat, err := wallet.GetTokenBalance(token, network)
	if err != nil {
		return "", fmt.Errorf("error getting account: %v", err)
	}

	balanceStr := strconv.FormatFloat(balanceFloat, 'f', 4, 64)

	return balanceStr, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var wallet *hdwallet.Wallet = nil
	provider := "hardhat"

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

			err = wallet.Initialize(password)
			if err != nil {
				fmt.Println("Error initializing wallet:", err)
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

			err = wallet.Initialize(password)
			if err != nil {
				fmt.Println("Error initializing wallet:", err)
				break
			}

			fmt.Println("Wallet restored successfully")

		case "get-token-balance":
			if wallet == nil {
				fmt.Println("Wallet not found")
				break
			}

			if !eth.NetListening(provider) {
				fmt.Println("Node is not listening")
				break
			}
			fmt.Println("Enter token name: ")
			scanner.Scan()
			token := strings.TrimSpace(scanner.Text())
			balance, err := GetBalance(wallet, token, provider)
			if err != nil {
				fmt.Println("Failed to get balance: ", err)
				break
			}

			fmt.Printf("Balance for %s token: %s\n", token, balance)

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
