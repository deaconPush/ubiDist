package main

import (
	"context"
	"fmt"
	"log"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx    context.Context
	wallet *utils.Wallet
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
	wallet, mnemonic, err := utils.CreateWallet(password)
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
	wallet, err := utils.RestoreWallet(password, mnemonic)
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

func (a *App) GetBalance() map[string]float64 {
	assets := map[string]float64{
		"ETH": 0,
	}
	var err error

	assets["ETH"], err = a.wallet.GetTokenBalance("ETH", "hardhat")

	if err != nil {
		log.Printf("error getting balance: %v", err)
	}

	fmt.Printf("Assets: %v", assets)

	return assets
}
