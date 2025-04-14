package hdwallet

import (
	"context"
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

	storeMasterKey(ws, password, masterKey)
	return &Wallet{publicKey: masterKey.PublicKey()}, mnemonic, nil
}

func RestoreWallet(password string, mnemonic string, ws *WalletStorage) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("error recovering master key from seed: %v", err)
	}

	err = storeMasterKey(ws, password, masterKey)
	if err != nil {
		return nil, fmt.Errorf("error storing master key: %v", err)
	}

	return &Wallet{publicKey: masterKey.PublicKey(), walletDB: ws}, nil
}

func RecoverWallet(password string) (*Wallet, error) {
	ws, err := utils.NewDatabaseService(context.Background(), "wallet.db")
	if err != nil {
		return nil, fmt.Errorf("error initializing database service: %v", err)
	}

	walletDB := NewWalletStorage(ws.GetDB())

	pubKeyHex, encryptedRootKey, err := walletDB.RetrieveKeysFromDB(password)
	if err != nil {
		return nil, fmt.Errorf("error retrieving keys from DB: %v", err)
	}

	pubKeyData, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding public key: %v", err)
	}

	_, err = utils.Decrypt([]byte(password), []byte(encryptedRootKey))
	if err != nil {
		return nil, fmt.Errorf("error decrypting master key %v", err)
	}

	pubKey, err := bip32.Deserialize(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing public key: %v", err)
	}

	wallet := &Wallet{
		publicKey: pubKey,
		walletDB:  walletDB,
	}

	return wallet, nil
}

func (w *Wallet) retrieveMasterKey(password string) (*bip32.Key, error) {
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing master public key: %v", err)
	}
	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(password, pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error retrieving key from DB: %v", err)
	}

	return masterKey, nil
}

func storeMasterKey(ws *WalletStorage, password string, masterKey *bip32.Key) error {
	masterKeyData, err := masterKey.Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master Key: %v", err)
	}

	pubKeyData, err := masterKey.PublicKey().Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master public key: %v", err)
	}

	masterKeyHex := hex.EncodeToString(masterKeyData)
	encryptedMasterKey, err := utils.Encrypt([]byte(password), []byte(masterKeyHex))
	if err != nil {
		return fmt.Errorf("error encrypting data: %v", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	err = ws.SaveRootKeyToDB(password, pubKeyHex, encryptedMasterKey)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %v", err)
	}

	return nil
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
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing master public key: %v", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(password, pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error retrieving key from DB: %v", err)
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
