package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"wallet/internal/currencies/eth"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	publicKey *bip32.Key
}

func (w *Wallet) retrieveRootKey() (*bip32.Key, error) {
	rootFile, err := os.Open("wallet.json")

	if err != nil {
		return nil, fmt.Errorf("error opening wallet file: %v", err)
	}

	defer rootFile.Close()

	keyBytes, err := io.ReadAll(rootFile)
	if err != nil {
		return nil, fmt.Errorf("error reading wallet file: %v", err)
	}

	var result map[string]interface{}
	json.Unmarshal(keyBytes, &result)
	keyHex, ok := result["HDKey"].(string)
	if !ok {
		return nil, fmt.Errorf("error reading HDKey from wallet file")
	}

	keyData, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding HDKey from wallet file: %v", err)
	}

	masterKey, err := bip32.Deserialize(keyData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing HDKey from wallet file: %v", err)
	}

	return masterKey, nil
}

func saveHDKey(key *bip32.Key) error {
	keyData, err := key.Serialize()
	filename := "wallet.json"
	if err != nil {
		return fmt.Errorf("error serializing HDKey: %v", err)
	}
	var data = map[string]interface{}{
		"HDKey": hex.EncodeToString(keyData),
	}
	keyBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling HDKey data: %v", err)
	}

	err = os.WriteFile(filename, keyBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing HDKey to file: %v", err)
	}
	return nil
}

func CreateWallet(password string) (*Wallet, string, error) {
	mnemonic, err := GenerateMnemonic()
	if err != nil {
		return nil, "", fmt.Errorf("error generating mnemonic: %v", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, "", fmt.Errorf("error recovering master key from seed:: %v", err)
	}

	err = saveHDKey(masterKey)
	if err != nil {
		return nil, "", fmt.Errorf("error saving HDKey: %v", err)
	}

	pubKey := masterKey.PublicKey()
	err = saveHDKey(masterKey)
	if err != nil {
		return nil, "", fmt.Errorf("error saving HDKey: %v", err)
	}

	return &Wallet{publicKey: pubKey}, mnemonic, nil
}

func RestoreWallet(password string, mnemonic string) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("error recovering master key from seed: %v", err)
	}

	err = saveHDKey(masterKey)
	if err != nil {
		return nil, fmt.Errorf("error saving HDKey: %v", err)
	}

	pubKey := masterKey.PublicKey()
	err = saveHDKey(masterKey)
	if err != nil {
		return nil, fmt.Errorf("error saving HDKey: %v", err)
	}

	return &Wallet{publicKey: pubKey}, nil
}

func (w *Wallet) CreateETHAccount() (*eth.ETHAccount, error) {
	masterKey, err := w.retrieveRootKey()
	if err != nil {
		return nil, fmt.Errorf("error retrieving root key: %v", err)
	}

	ethKey, err := DeriveChildKey(masterKey, "m/44'/60'/0'/0/0")
	if err != nil {
		return nil, fmt.Errorf("error deriving ETH key: %v", err)
	}

	ethAccount, err := eth.NewETHAccount(ethKey, "ETH")
	if err != nil {
		return nil, fmt.Errorf("error creating ETH account: %v", err)
	}
	return ethAccount, nil
}
