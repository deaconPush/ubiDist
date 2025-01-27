package main

import (
	"context"

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
	if !bip39.IsMnemonicValid(mnemonic) {
		return false
	}
	return true
}
