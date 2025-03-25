package hdwallet

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strconv"
	"wallet/internal/currencies/eth"
	"wallet/internal/utils"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	publicKey *bip32.Key
	Accounts  []Account
	walletDB  *WalletStorage
}

type Account interface {
	GetAddress() string
	RetrieveBalance(network string) (string, error)
	GetTokenName() string
}

func (w *Wallet) retrieveRootKey(password string) (*bip32.Key, error) {
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing master public key: %v", err)
	}
	pubKeyHex := hex.EncodeToString(pubKeyData)
	var encryptedKeyData string
	err = w.db.QueryRowContext(context.Background(), "SELECT masterKey FROM wallets where publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
	switch {
	case err == sql.ErrNoRows:
		return nil, fmt.Errorf("no rows returned")
	case err != nil:
		return nil, fmt.Errorf("error querying database: %v", err)
	}

	keyDataHex, err := utils.Decrypt([]byte(password), []byte(encryptedKeyData))
	if err != nil {
		return nil, fmt.Errorf("error decrypting master key %v", err)
	}

	keyData, err := hex.DecodeString(string(keyDataHex))
	if err != nil {
		return nil, fmt.Errorf("error decoding key data: %v", err)
	}

	masterKey, err := bip32.Deserialize(keyData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing HDKey from wallet file: %v", err)
	}

	return masterKey, nil
}

func StoreMasterKey(ws *WalletStorage, password string, masterKey *bip32.Key) error {
	masterKeyData, err := masterKey.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master Key: %v", err)
	}

	pubKeyData, err := masterKey.PublicKey().Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master public key: %v", err)
	}

	masterKeyHex := hex.EncodeToString(masterKeyData)
	encryptedMasterKeyHex, err := utils.Encrypt([]byte(password), []byte(masterKeyHex))
	if err != nil {
		return fmt.Errorf("error encrypting data: %v", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	err = ws.SaveRootKeyToDB(password, pubKeyHex, encryptedMasterKeyHex)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %v", err)
	}

	return nil
}

func CreateWallet(password string, ws *WalletStorage) (*Wallet, string, error) {
	mnemonic, err := utils.GenerateMnemonic()
	if err != nil {
		return nil, "", fmt.Errorf("error generating mnemonic: %v", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, "", fmt.Errorf("error recovering master key from seed:: %v", err)
	}

	StoreMasterKey(ws, password, masterKey)
	return &Wallet{publicKey: masterKey.PublicKey()}, mnemonic, nil
}

func RestoreWallet(password string, mnemonic string, ws *WalletStorage) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("error recovering master key from seed: %v", err)
	}

	StoreMasterKey(ws, password, masterKey)
	return &Wallet{publicKey: masterKey.PublicKey(), walletDB: ws}, nil
}

func (w *Wallet) Initialize(password string) error {
	ethAccount, err := w.CreateETHAccount(password)
	if err != nil {
		return fmt.Errorf("error creating ETH account: %v", err)
	}

	w.Accounts = append(w.Accounts, ethAccount)
	return nil
}

func (w *Wallet) CreateETHAccount(password string) (*eth.ETHAccount, error) {
	masterKey, err := w.retrieveRootKey(password)
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

func (w *Wallet) GetTokenBalance(tokenName string, options ...string) (float64, error) {
	var network string = ""
	if len(options) > 0 {
		network = options[0]
	}
	for _, account := range w.Accounts {
		if account.GetTokenName() == tokenName {
			hexBalance, err := account.RetrieveBalance(network)
			if err != nil {
				return 0, fmt.Errorf("error retrieving balance: %v", err)
			}

			balance, err := eth.HexToEther(hexBalance)
			if err != nil {
				return 0, fmt.Errorf("error converting balance: %v", err)
			}

			return strconv.ParseFloat(balance, 64)
		}
	}
	return 0, fmt.Errorf("token not found: %s", tokenName)
}

func (w *Wallet) GetAccounts() []Account {
	return w.Accounts

}
