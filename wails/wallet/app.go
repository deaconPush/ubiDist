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

type RPCPayload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
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

func getGasPrice(provider string) string {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		ID:      uuid.New().String(),
	}

	response := sendRPCRequest(payload, provider)
	gasPriceHex := response["result"].(string)
	return gasPriceHex
}

func getTransactionCount(provider string, address string) string {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionCount",
		Params:  []interface{}{address, "latest"},
		ID:      uuid.New().String(),
	}

	response := sendRPCRequest(payload, provider)
	transactionCountHex := response["result"].(string)
	return transactionCountHex
}

func sendRPCRequest(payload RPCPayload, url string) map[string]interface{} {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}

	respBytes, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}

	defer respBytes.Body.Close()

	body, err := io.ReadAll(respBytes.Body)
	if err != nil {
		log.Fatal(err)
	}

	response := map[string]interface{}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func getBalance(provider string, address string) string {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{address, "latest"},
		ID:      uuid.New().String(),
	}

	response := sendRPCRequest(payload, provider)
	balanceHex := response["result"].(string)
	return balanceHex
}

func weiToEther(wei string) string {
	weiBigInt := new(big.Int)
	weiBigInt.SetString(wei, 0)
	ether := new(big.Float).SetInt(weiBigInt)
	ether = new(big.Float).Quo(ether, big.NewFloat(1e18))
	return ether.String()
}
