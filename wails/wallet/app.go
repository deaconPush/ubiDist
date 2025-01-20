package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
)

// App struct
type App struct {
	ctx context.Context
}

type Transaction struct {
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	To       string
	Value    *big.Int
	Data     []byte
	V        *big.Int
	R        *big.Int
	S        *big.Int
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

func getGasPrice(provider string) *big.Int {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		ID:      uuid.New().String(),
	}

	response := sendRPCRequest(payload, provider)
	gasPriceHex := response["result"].(string)
	gasPrice, _ := new(big.Int).SetString(gasPriceHex[2:], 16)
	return gasPrice
}

func getNonce(provider string, address string) uint64 {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionCount",
		Params:  []interface{}{address, "latest"},
		ID:      uuid.New().String(),
	}

	response := sendRPCRequest(payload, provider)
	nonceHex := response["result"].(string)
	nonce := new(big.Int)
	nonce.SetString(nonceHex[2:], 16)
	return nonce.Uint64()
}

func estimateGas(provider string, from string, to string, value *big.Int) uint64 {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_estimateGas",
		Params: []interface{}{
			map[string]interface{}{
				"from":  from,
				"to":    to,
				"value": value,
			},
		},
		ID: uuid.New().String(),
	}

	response := sendRPCRequest(payload, provider)
	gasLimitHex := response["result"].(string)
	gasLimit := new(big.Int)
	gasLimit.SetString(gasLimitHex[2:], 16)
	return gasLimit.Uint64()
}

func signTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey, chainID *big.Int) ([]byte, error) {
	txRLP := []interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		chainID, big.NewInt(0), big.NewInt(0),
	}

	encodedTx, err := rlp.EncodeToBytes(txRLP)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to RLP encode transaction: %w", err)
	}

	txHash := crypto.Keccak256(encodedTx)
	signature, err := crypto.Sign(txHash, privateKey)
	fmt.Println("signature:", signature)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	r := new(big.Int).SetBytes(signature[0:32])
	s := new(big.Int).SetBytes(signature[32:64])
	v := big.NewInt(int64(signature[64]) + 35 + 2*chainID.Int64())
	// Add the signature to the transaction
	tx.V = v
	tx.R = r
	tx.S = s
	// RPL encode the transaction with the signature
	signedTxRLP := []interface{}{
		tx.Nonce,
		tx.GasPrice,
		tx.GasLimit,
		tx.To,
		tx.Value,
		tx.Data,
		tx.V,
		tx.R,
		tx.S,
	}
	fmt.Println("Signing transaction with appropiate data...")
	signedEncodedTx, err := rlp.EncodeToBytes(signedTxRLP)
	fmt.Println("signedEncodedTx:", signedEncodedTx)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to RLP encode signed transaction: %w", err)
	}

	return signedEncodedTx, nil
}

func processTransaction(provider string, from string, to string, value *big.Int, privateKey *ecdsa.PrivateKey, chainID int64) string {
	nonce := getNonce(provider, from)
	gasPrice := getGasPrice(provider)
	gasLimit := estimateGas(provider, from, to, value)
	tx := Transaction{
		Nonce:    nonce,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		To:       to,
		Value:    value,
		Data:     nil,
	}

	fmt.Print("Signing transaction...")
	signedTx, err := signTransaction(&tx, privateKey, big.NewInt(chainID))
	if err != nil {
		log.Fatal("Error signing transaction :", err)
	}

	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{signedTx},
		ID:      uuid.New().String(),
	}
	if err != nil {
		log.Fatal(err)
	}

	response := sendRPCRequest(payload, provider)
	fmt.Println("response:", response)
	return ""
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
