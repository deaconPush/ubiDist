package eth

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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
)

type Transaction struct {
	Nonce    uint64
	GasPrice *big.Int
	GasLimit uint64
	To       common.Address
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

	response := map[string]interface{}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return response, nil
}

func NetListening(provider string) bool {
	payload := RPCPayload{
		Jsonrpc: "2.0",
		Method:  "net_listening",
		Params:  []interface{}{},
		ID:      uuid.New().String(),
	}

	response, err := sendRPCRequest(payload, provider)
	if err != nil {
		return false
	}

	return response["result"].(bool)
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

func ProcessTransactionWithNativeSigning(provider string, from string, to string, value *big.Int, privateKey *ecdsa.PrivateKey) (string, error) {
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
	tx := Transaction{
		Nonce:    nonce,
		GasPrice: gasPrice, // 20 Gwei
		GasLimit: gasLimit, // Standard ETH transfer
		To:       toAddress,
		Value:    value,
		Data:     nil,
	}
	signedTx, err := signTransaction(&tx, privateKey, big.NewInt(chainID))
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	rawTxHex := hex.EncodeToString(signedTx)

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
	signedEncodedTx, err := rlp.EncodeToBytes(signedTxRLP)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to RLP encode signed transaction: %w", err)
	}

	return signedEncodedTx, nil
}

func HexToEther(hexBalance string) (string, error) {
	balance, ok := new(big.Int).SetString(hexBalance[2:], 16)
	if !ok {
		return "", fmt.Errorf("failed to convert hex to big.Int")
	}
	ether := new(big.Float).SetInt(balance)
	// Convert wei to ether
	ether.Quo(ether, big.NewFloat(1e18))

	return ether.String(), nil
}

func EtherToWei(ether string) (*big.Int, error) {
	etherFloat, _, err := new(big.Float).Parse(ether, 10)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ether: %w", err)
	}

	wei := new(big.Int)
	etherInt := new(big.Int)
	etherFloat.Mul(etherFloat, big.NewFloat(1e18))
	etherFloat.Int(etherInt)
	wei.Set(etherInt)

	return wei, nil
}
