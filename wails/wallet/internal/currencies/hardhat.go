package currencies

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
)

type RPCPayload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
}

func sendRPCRequest(payload RPCPayload, url string) (map[string]interface{}, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	respBytes, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer respBytes.Body.Close()

	body, err := io.ReadAll(respBytes.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("response body:", string(body))
	response := map[string]interface{}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return response, nil
}

func GetGasPrice(provider string) (*big.Int, error) {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_gasPrice",
		Params:  []interface{}{},
		ID:      uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get gas price: %w", err)
	}
	gasPriceHex := response["result"].(string)
	gasPrice, _ := new(big.Int).SetString(gasPriceHex[2:], 16)
	return gasPrice, nil
}

func GetNonce(provider string, address string) (uint64, error) {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_getTransactionCount",
		Params:  []interface{}{address, "latest"},
		ID:      uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return 0, fmt.Errorf("failed to get nonce: %w", err)
	}
	nonceHex := response["result"].(string)
	nonce := new(big.Int)
	nonce.SetString(nonceHex[2:], 16)
	return nonce.Uint64(), nil
}

func EstimateGas(provider string, from string, to string, value *big.Int) (uint64, error) {
	valueHex := fmt.Sprintf("0x%x", value)
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_estimateGas",
		Params: []interface{}{
			map[string]interface{}{
				"from":  from,
				"to":    to,
				"value": valueHex,
			},
		},
		ID: uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}
	gasLimitHex := response["result"].(string)
	gasLimit := new(big.Int)
	gasLimit.SetString(gasLimitHex[2:], 16)
	return gasLimit.Uint64(), nil
}

func GetChainId(provider string) (int64, error) {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_chainId",
		Params:  []interface{}{},
		ID:      uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return -1, fmt.Errorf("failed to get chain ID: %w", err)
	}

	chainIdHex := response["result"].(string)
	// Convert hex to int64
	chainId, _ := new(big.Int).SetString(chainIdHex[2:], 16)
	return chainId.Int64(), nil
}

func ProcessTransaction(provider string, from string, to string, value *big.Int, privateKey *ecdsa.PrivateKey) (string, error) {
	toAddress := common.HexToAddress(to)
	nonce, err := GetNonce(provider, from)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve nonce: %w", err)
	}

	gasPrice, err := GetGasPrice(provider)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve gas price: %w", err)
	}

	gasLimit, err := EstimateGas(provider, from, to, value)
	if err != nil {
		return "", fmt.Errorf("failed to estimate gas: %w", err)
	}

	chainID, err := GetChainId(provider)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve chain ID: %w", err)
	}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainID)), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	rawTx := types.Transactions{signedTx}
	rawTxBytes, err := rlp.EncodeToBytes(rawTx[0])
	if err != nil {
		return "", fmt.Errorf("failed to RLP encode transaction: %w", err)
	}

	rawTxHex := hex.EncodeToString(rawTxBytes)
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []interface{}{rawTxHex},
		ID:      uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	txHash := response["result"].(string)
	return txHash, nil
}

func GetBalance(provider string, address string) (string, error) {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "eth_getBalance",
		Params:  []interface{}{address, "latest"},
		ID:      uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return "", fmt.Errorf("failed to get balance: %w", err)
	}
	balanceHex := response["result"].(string)
	return balanceHex, nil
}
