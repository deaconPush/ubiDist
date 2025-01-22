package main

import (
	"context"
	"wallet/internal/currencies"
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

func (a *App) GetBalance() (string, error) {
	address := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	hexBalance, err := currencies.GetBalance("http://localhost:8545", address)
	if err != nil {
		return "", err
	}

	etherBalance, err := currencies.HexToEther(hexBalance)
	if err != nil {
		return "", err
	}

	return etherBalance, nil
}
