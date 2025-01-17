package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"

	"github.com/google/uuid"
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

func getBalance(provider string, address string) string {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBalance",
		"params":  []interface{}{address, "latest"},
		"id":      uuid.New().String(),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	respBytes, err := http.Post(provider, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
		return ""
	}

	defer respBytes.Body.Close()

	body, err := io.ReadAll(respBytes.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	response := map[string]interface{}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	balanceHex := response["result"].(string)
	balance := weiToEther(balanceHex)
	return balance
}

func weiToEther(wei string) string {
	weiBigInt := new(big.Int)
	weiBigInt.SetString(wei, 0)
	ether := new(big.Float).SetInt(weiBigInt)
	ether = new(big.Float).Quo(ether, big.NewFloat(1e18))
	return ether.String()
}
