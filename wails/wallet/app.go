package main

import (
	"context"
	"fmt"
	"wallet/internal/hdwallet"

	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx    context.Context
	wallet *hdwallet.Wallet
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func (a *App) CreateWallet(password string) (string, error) {
	wallet, mnemonic, err := hdwallet.CreateWallet(password)
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
	wallet, err := hdwallet.RestoreWallet(password, mnemonic)
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

func (a *App) GetAssets(tokens []string) map[string]float64 {
	var assets = make(map[string]float64)
	for _, token := range tokens {
		balance, err := a.wallet.GetTokenBalance(token)
		if err != nil {
			fmt.Println("error getting balance for token:", err)
		}

		assets[token] = balance
	}
	return assets
}
