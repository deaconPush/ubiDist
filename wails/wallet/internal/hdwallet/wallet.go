package hdwallet

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
	"wallet/internal/currencies/eth"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	publicKey *bip32.Key
	Accounts  []Account
	walletDB  *WalletStorage
	ctx       context.Context
}

type Account interface {
	GetAddress() string
	RetrieveBalance() (string, error)
	GetTokenName() string
	EstimateGas(from, value string) (string, error)
}

func CreateWallet(ctx context.Context, password string, ws *WalletStorage) (*Wallet, string, error) {
	mnemonic, err := utils.GenerateMnemonic()
	if err != nil {
		return nil, "", fmt.Errorf("error generating mnemonic: %v", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, "", fmt.Errorf("error recovering master key from seed:: %v", err)
	}

	storeMasterKey(ctx, ws, password, masterKey)
	return &Wallet{publicKey: masterKey.PublicKey(), walletDB: ws, ctx: ctx}, mnemonic, nil
}

func RestoreWallet(ctx context.Context, password string, mnemonic string, ws *WalletStorage) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("error recovering master key from seed: %v", err)
	}

	err = storeMasterKey(ctx, ws, password, masterKey)
	if err != nil {
		return nil, fmt.Errorf("error storing master key: %v", err)
	}

	return &Wallet{publicKey: masterKey.PublicKey(), walletDB: ws, ctx: ctx}, nil
}

func RecoverWallet(ctx context.Context, password string, ws *WalletStorage) (*Wallet, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	pubKey, err := ws.RetrievePublicKeyFromDB(dbCtx)
	defer cancel()
	if err != nil {
		return nil, fmt.Errorf("error retrieving public key from DB: %v", err)
	}

	wallet := &Wallet{
		publicKey: pubKey,
		walletDB:  ws,
		ctx:       ctx,
	}
	return wallet, nil
}

func storeMasterKey(ctx context.Context, ws *WalletStorage, password string, masterKey *bip32.Key) error {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
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
	err = ws.SaveRootKeyToDB(dbCtx, password, pubKeyHex, encryptedMasterKey)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %v", err)
	}

	return nil
}

func (w *Wallet) Initialize(password string) error {
	dbCtx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()
	ethAccount, err := w.CreateETHAccount(dbCtx, password)
	if err != nil {
		return fmt.Errorf("error creating ETH account: %v", err)
	}

	w.Accounts = append(w.Accounts, ethAccount)
	return nil
}

func (w *Wallet) CreateETHAccount(ctx context.Context, password string) (*eth.ETHAccount, error) {
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing master public key: %v", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(ctx, password, pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error retrieving key from DB: %v", err)
	}

	ethKey, err := utils.DeriveChildKey(masterKey, "m/44'/60'/0'/0/0")
	if err != nil {
		return nil, fmt.Errorf("error deriving ETH key: %v", err)
	}

	ethAccount, err := eth.NewETHAccount(w.ctx, ethKey, "ETH")
	if err != nil {
		return nil, fmt.Errorf("error creating ETH account: %v", err)
	}
	return ethAccount, nil
}

func (w *Wallet) GetTokenBalance(tokenName string) (float64, error) {
	for _, account := range w.Accounts {
		if account.GetTokenName() == tokenName {
			hexBalance, err := account.RetrieveBalance()
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

func (w *Wallet) EstimateTokenGas(tokenName, to, value string) (string, error) {
	for _, account := range w.Accounts {
		if account.GetTokenName() == tokenName {
			gasPrice, err := account.EstimateGas(to, value)
			if err != nil {
				return "", fmt.Errorf("Error estimating gad price for token %s : %v", tokenName, err)
			}

			return gasPrice, nil
		}
	}
	return "", fmt.Errorf("token not found: %s", tokenName)
}
