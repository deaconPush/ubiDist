package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"wallet/internal/currencies/eth"
	"wallet/internal/hdwallet"
)

func createWallet(ctx context.Context, password string) (*hdwallet.Wallet, error) {
	walletDB, err := hdwallet.NewWalletStorage(ctx, ":memory:")
	if err != nil {
		panic(fmt.Errorf("error initializing wallet storage: %w", err))
	}

	wallet, _, err := hdwallet.CreateWallet(ctx, password, walletDB)
	if err != nil {
		return nil, fmt.Errorf("error creating wallet: %w", err)
	}

	return wallet, nil
}

func RestoreWallet(ctx context.Context, password, mnemonic string) (*hdwallet.Wallet, error) {
	walletDB, err := hdwallet.NewWalletStorage(ctx, ":memory:")
	if err != nil {
		panic(fmt.Errorf("error initializing wallet storage: %w", err))
	}

	wallet, err := hdwallet.RestoreWallet(ctx, password, mnemonic, walletDB)
	if err != nil {
		return nil, fmt.Errorf("error restoring wallet: %w", err)
	}

	return wallet, nil
}

func GetBalance(wallet *hdwallet.Wallet, token string) (string, error) {
	balanceFloat, err := wallet.GetBalance(token, 0)
	if err != nil {
		return "", fmt.Errorf("error getting account: %w", err)
	}

	balanceStr := strconv.FormatFloat(balanceFloat, 'f', 4, 64)

	return balanceStr, nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	tokens := []string{"ETH"}
	var wallet *hdwallet.Wallet = nil
	cliCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client := eth.NewEthClient("http://localhost:8545")
	defer cancel()

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
			wallet, err = createWallet(context.Background(), password)
			if err != nil {
				fmt.Println("Error creating wallet:", err)
				break
			}

			err = wallet.Initialize(tokens, password)
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
			wallet, err = RestoreWallet(context.Background(), password, mnemonic)
			if err != nil {
				fmt.Println("Error restoring wallet:", err)
				break
			}

			err = wallet.Initialize(tokens, password)
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

			if !client.NetListening(cliCtx) {
				fmt.Println("Node is not listening")
				break
			}
			fmt.Println("Enter token name: ")
			scanner.Scan()
			token := strings.TrimSpace(scanner.Text())
			balance, err := GetBalance(wallet, token)
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
