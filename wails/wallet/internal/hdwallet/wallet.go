package hdwallet

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"wallet/internal/currencies/eth"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	publicKey *bip32.Key
	Accounts  map[string]Account
}

type Account interface {
	GetAddress() string
	RetrieveBalance(network string) (string, error)
	GetTokenName() string
}

func (w *Wallet) retrieveRootKey() (*bip32.Key, error) {
	rootFile, err := os.Open("wallet.json")
	if err != nil {
		return nil, fmt.Errorf("error opening wallet file: %v", err)
	}

	defer rootFile.Close()

	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing public key: %v", err)
	}

	pubKeyHexData := hex.EncodeToString(pubKeyData)
	keyBytes, err := io.ReadAll(rootFile)
	if err != nil {
		return nil, fmt.Errorf("error reading wallet file: %v", err)
	}

	result := map[string]string{}
	json.Unmarshal(keyBytes, &result)
	keyData, ok := result[pubKeyHexData]
	if !ok {
		return nil, fmt.Errorf("error retrieving HDKey from wallet file")
	}

	keyDataBytes, err := hex.DecodeString(keyData)
	if err != nil {
		return nil, fmt.Errorf("error decoding key data: %v", err)
	}

	masterKey, err := bip32.Deserialize(keyDataBytes)
	if err != nil {
		return nil, fmt.Errorf("error deserializing HDKey from wallet file: %v", err)
	}

	return masterKey, nil
}

func saveHDKey(masterKey, pubKey *bip32.Key) error {
	masterKeyData, err := masterKey.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master Key: %v", err)
	}

	pubKeyData, err := pubKey.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master public key: %v", err)
	}

	masterKeyHexData := hex.EncodeToString(masterKeyData)
	pubKeyHexData := hex.EncodeToString(pubKeyData)
	filename := "wallet.json"

	if err != nil {
		return fmt.Errorf("error serializing HDKey: %v", err)
	}

	var data = map[string]interface{}{
		pubKeyHexData: masterKeyHexData,
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
	mnemonic, err := utils.GenerateMnemonic()
	if err != nil {
		return nil, "", fmt.Errorf("error generating mnemonic: %v", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, "", fmt.Errorf("error recovering master key from seed:: %v", err)
	}
	pubKey := masterKey.PublicKey()
	err = saveHDKey(masterKey, pubKey)
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

	pubKey := masterKey.PublicKey()
	err = saveHDKey(masterKey, pubKey)
	if err != nil {
		return nil, fmt.Errorf("error saving HDKey: %v", err)
	}

	return &Wallet{publicKey: pubKey}, nil
}

func (w *Wallet) Initialize(password string) error {
	var err error
	w.Accounts = make(map[string]Account)
	w.Accounts["ETH"], err = w.CreateETHAccount(password)
	if err != nil {
		return fmt.Errorf("error creating ETH account: %v", err)
	}

	return nil
}

func (w *Wallet) CreateETHAccount(password string) (*eth.ETHAccount, error) {
	masterKey, err := w.retrieveRootKey()
	if err != nil {
		return nil, fmt.Errorf("error retrieving root key: %v", err)
	}

	ethKey, err := utils.DeriveChildKey(masterKey, "m/44'/60'/0'/0/0")
	if err != nil {
		return nil, fmt.Errorf("error deriving ETH key: %v", err)
	}

	ethAccount, err := eth.NewETHAccount(ethKey, "ETH")
	if err != nil {
		return nil, fmt.Errorf("error creating ETH account: %v", err)
	}
	return ethAccount, nil
}

func (w *Wallet) GetTokenBalance(tokenName, network string) (float64, error) {
	for _, account := range w.Accounts {
		if account.GetTokenName() == tokenName {
			balance, err := account.RetrieveBalance(network)
			if err != nil {
				return 0, fmt.Errorf("error retrieving balance: %v", err)
			}
			return strconv.ParseFloat(balance, 64)
		}
	}
	return 0, fmt.Errorf("token not found: %s", tokenName)
}
