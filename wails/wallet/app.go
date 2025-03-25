package main

import (
	"context"
	"fmt"
	"time"
	"wallet/internal/hdwallet"
	"wallet/internal/utils"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx       context.Context
	wallet    *hdwallet.Wallet
	dbService *utils.DatabaseService
	walletDB  *hdwallet.WalletStorage
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbService, err := utils.NewDatabaseService(dbCtx)
	if err != nil {
		panic(fmt.Errorf("error initializing database service: %v", err))
	}

	walletDB := hdwallet.NewWalletStorage(dbService.GetDB())
	a.walletDB = walletDB
	a.dbService = dbService
}

func (a *App) WalletExists() (bool, error) {
	exists, err := a.walletDB.WalletExists(a.ctx)
	if err != nil {
		return false, fmt.Errorf("error checking if wallet exists: %v", err)
	}

	return exists, nil
}

func (a *App) ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func (a *App) CreateWallet(password string) (string, error) {
	wallet, mnemonic, err := hdwallet.CreateWallet(password, a.walletDB)
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
	wallet, err := hdwallet.RestoreWallet(password, mnemonic, a.walletDB)
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

func (a *App) GetAssets(tokenSymbols []string) map[string]float64 {
	var assets = make(map[string]float64)
	for _, token := range tokenSymbols {
		balance, err := a.wallet.GetTokenBalance(token)
		if err != nil {
			fmt.Println("error getting balance for token:", err)
		}

		assets[token] = balance
	}
	return assets
}
