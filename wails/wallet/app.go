package main

import (
	"context"
	"database/sql"
	"fmt"
	"wallet/internal/hdwallet"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx    context.Context
	wallet *hdwallet.Wallet
	db     *sql.DB
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "wallet.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	return db, nil
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	db, err := initDB()
	if err != nil {
		panic(fmt.Sprintf("error initializing database: %v", err))
	}

	_, err = db.ExecContext(a.ctx, "CREATE TABLE IF NOT EXISTS wallets (publicKey TEXT PRIMARY KEY, masterKey TEXT)")
	if err != nil {
		panic(fmt.Sprintf("error creating accounts table: %v", err))
	}

	a.db = db
}

func (a *App) WalletExists() (bool, error) {
	var count int
	db := a.db
	err := db.QueryRowContext(a.ctx, "SELECT COUNT(*) FROM wallets").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if wallet exists: %v", err)
	}

	return count > 0, nil
}

func (a *App) ValidateMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

func (a *App) CreateWallet(password string) (string, error) {
	wallet, mnemonic, err := hdwallet.CreateWallet(password, a.db)
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
	wallet, err := hdwallet.RestoreWallet(password, mnemonic, a.db)
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
