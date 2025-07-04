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

func createWalletCmd(scanner *bufio.Scanner, tokens []string) (*hdwallet.Wallet, error) {
	fmt.Fprintln(os.Stdout, "Enter password: ")
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read password input")
	}
	password := strings.TrimSpace(scanner.Text())

	wallet, err := createWallet(context.Background(), password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating wallet:", err)
		return nil, err
	}

	err = wallet.Initialize(tokens, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing wallet:", err)
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "Wallet created successfully.")
	return wallet, nil
}

func restoreWalletCmd(scanner *bufio.Scanner, tokens []string) (*hdwallet.Wallet, error) {
	fmt.Fprintln(os.Stdout, "Enter password: ")
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read password")
	}
	password := strings.TrimSpace(scanner.Text())

	fmt.Fprintln(os.Stdout, "Enter mnemonic: ")
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read mnemonic")
	}
	mnemonic := strings.TrimSpace(scanner.Text())

	fmt.Fprintln(os.Stdout, mnemonic)

	wallet, err := RestoreWallet(context.Background(), password, mnemonic)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error restoring wallet:", err)
		return nil, err
	}

	err = wallet.Initialize(tokens, password)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing wallet:", err)
		return nil, err
	}

	fmt.Fprintln(os.Stdout, "Wallet restored successfully")
	return wallet, nil
}

func checkBalanceCmd(ctx context.Context, scanner *bufio.Scanner, wallet *hdwallet.Wallet, client *eth.Client) error {
	if wallet == nil {
		fmt.Fprintln(os.Stderr, "Wallet not found")
		return fmt.Errorf("wallet not initialized")
	}

	if !client.NetListening(ctx) {
		fmt.Fprintln(os.Stderr, "Node is not listening")
		return fmt.Errorf("node not listening")
	}

	fmt.Fprintln(os.Stdout, "Enter token name: ")
	if !scanner.Scan() {
		return fmt.Errorf("failed to read token input")
	}
	token := strings.TrimSpace(scanner.Text())

	balance, err := GetBalance(wallet, token)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get balance:", err)
		return err
	}

	fmt.Fprintf(os.Stdout, "Balance for %s token: %s\n", token, balance)
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	tokens := []string{"ETH"}
	var wallet *hdwallet.Wallet
	var err error
	cliCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	client := eth.NewClient("http://localhost:8545")
	defer cancel()

	for {
		fmt.Fprintln(os.Stdout, "Enter a command or (type 'exit' to quit):")
		scanner.Scan()
		command := scanner.Text()
		switch command {
		case "create-wallet":
			wallet, err = createWalletCmd(scanner, tokens)
			if err != nil {
				break
			}
		case "restore-wallet":
			wallet, err = restoreWalletCmd(scanner, tokens)
			if err != nil {
				break
			}
		case "get-token-balance":
			err := checkBalanceCmd(cliCtx, scanner, wallet, client)
			if err != nil {
				break
			}
		}
	}
}
