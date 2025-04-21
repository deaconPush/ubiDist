package main

import (
	"context"
	"fmt"
	"time"
	"wallet/internal/hdwallet"

	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx      context.Context
	wallet   *hdwallet.Wallet
	walletDB *hdwallet.WalletStorage
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	walletDB, err := hdwallet.NewWalletStorage("wallet.db", dbCtx)
	defer cancel()

	if err != nil {
		fmt.Println("error creating wallet storage:", err)
		return
	}

	a.walletDB = walletDB
}

func (a *App) WalletExists() (bool, error) {
	dbCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	exists, err := a.walletDB.WalletExists(dbCtx)
	defer cancel()
	if err != nil {
		return false, fmt.Errorf("error checking if wallet exists: %v", err)
	}

	return exists, nil
}

func (a *App) ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func (a *App) CreateWallet(password string) (string, error) {
	wallet, mnemonic, err := hdwallet.CreateWallet(a.ctx, password, a.walletDB)
	if err != nil {
		return "", fmt.Errorf("error creating wallet: %v", err)
	}

	a.wallet = wallet
	err = wallet.Initialize(password)
	if err != nil {
		return "", fmt.Errorf("error initializing wallet: %v", err)
	}

	return mnemonic, nil
}

func (a *App) RestoreWallet(password, mnemonic string) error {
	wallet, err := hdwallet.RestoreWallet(a.ctx, password, mnemonic, a.walletDB)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %v", err)
	}

	a.wallet = wallet
	err = wallet.Initialize(password)
	if err != nil {
		return fmt.Errorf("error initializing wallet: %v", err)
	}

	return nil
}

func (a *App) RecoverWallet(password string) error {
	wallet, err := hdwallet.RecoverWallet(a.ctx, password, a.walletDB)
	if err != nil {
		return fmt.Errorf("error recovering wallet: %v", err)
	}

	a.wallet = wallet
	err = wallet.Initialize(password)
	if err != nil {
		return fmt.Errorf("error initializing wallet: %v", err)
	}

	return nil
}

func (a *App) GetAssets(tokenSymbols []string) (map[string]float64, error) {
	var assets = make(map[string]float64)
	for _, token := range tokenSymbols {
		balance, err := a.wallet.GetTokenBalance(token)
		if err != nil {
			return nil, fmt.Errorf("error getting balance for token %s: %v", token, err)
		}

		assets[token] = balance
	}
	return assets, nil
}
