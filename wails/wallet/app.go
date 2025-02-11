package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"wallet/internal/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// App struct
type App struct {
	ctx    context.Context
	wallet *Wallet
}

type Wallet struct {
	Accounts []*Account
}

type Account struct {
	currency   string
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
}

func (w *Wallet) RetrieveRootKey() (*bip32.Key, error) {
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

func deriveChildKey(masterKey *bip32.Key, path string) (*bip32.Key, error) {
	indices, err := parseDerivationPath(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing derivation path: %v", err)
	}
	key := masterKey
	for _, index := range indices {
		key, err = key.NewChildKey(index)
		if err != nil {
			return nil, fmt.Errorf("error deriving child key: %v", err)
		}
	}
	return key, nil
}

func parseDerivationPath(path string) ([]uint32, error) {
	var indices []uint32
	var hardenedOffset uint32 = 0x80000000
	for _, part := range strings.Split(path, "/") {
		if part == "m" {
			continue
		}

		hardened := strings.HasSuffix(part, "'")
		if hardened {
			part = part[:len(part)-1]
		}

		parsedIndex, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing derivation path: %v", err)
		}
		index := uint32(parsedIndex)

		if hardened {
			index += hardenedOffset
		}

		indices = append(indices, index)
	}
	return indices, nil
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

	err = a.saveHDKey(masterKey)
	if err != nil {
		return "", fmt.Errorf("error saving HDKey: %v", err)
	}
	ethKey, err := deriveChildKey(masterKey, "m/44'/60'/0'/0/0")
	if err != nil {
		return "", fmt.Errorf("error deriving ETH key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(ethKey.Key)
	if err != nil {
		return "", fmt.Errorf("error deriving ETH private key: %v", err)
	}

	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	ethAccount := &Account{
		currency:   "ETH",
		privateKey: privateKey,
		publicKey:  publicKey,
	}
	wallet := &Wallet{
		Accounts: []*Account{ethAccount},
	}
	a.wallet = wallet

	return mnemonic, nil
}

func (a *App) RestoreWallet(password string, mnemonic string) error {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return fmt.Errorf("error recovering master key from seed: %v", err)
	}

	err = a.saveHDKey(masterKey)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %v", err)
	}
	return nil
}

func (a *App) saveHDKey(key *bip32.Key) error {
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
