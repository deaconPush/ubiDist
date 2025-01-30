package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx context.Context
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
	mnemonic, err := utils.GenerateMnemonic()
	if err != nil {
		return "", fmt.Errorf("error generating mnemonic: %v", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", fmt.Errorf("error recovering master key from seed:: %v", err)
	}

	fmt.Println("public key: ", hex.EncodeToString(masterKey.PublicKey().Key))
	return mnemonic, nil
}

func (a *App) RestoreWallet(password string, mnemonic string) error {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return fmt.Errorf("error recovering master key from seed: %v", err)
	}

	fmt.Println("public key: ", hex.EncodeToString(masterKey.PublicKey().Key))
	return nil
}
