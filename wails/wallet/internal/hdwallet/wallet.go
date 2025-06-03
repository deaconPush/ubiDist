package hdwallet

import (
	"context"
	"database/sql"
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
	Accounts  map[string]masterAccount
	walletDB  *WalletStorage
	ctx       context.Context
}

type masterAccount interface {
	GetAddress(accountIndex int) (string, error)
	RetrieveBalance(accountIndex int) (string, error)
	EstimateGas(from, value string, accountIndex int) (string, error)
	SendTransaction(to, value string, privateKey *bip32.Key, accountIndex int) (string, error)
}

type masterAccountFactory func(ctx context.Context, masterKey *bip32.Key, db *sql.DB) (masterAccount, error)

var masterAccountFactories = map[string]masterAccountFactory{
	"ETH": createETHAccount,
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
	return &Wallet{publicKey: masterKey.PublicKey(), Accounts: make(map[string]masterAccount), walletDB: ws, ctx: ctx}, mnemonic, nil
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

	return &Wallet{publicKey: masterKey.PublicKey(), Accounts: make(map[string]masterAccount), walletDB: ws, ctx: ctx}, nil
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
		Accounts:  make(map[string]masterAccount),
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

func (w *Wallet) Initialize(tokens []string, password string) error {
	for _, token := range tokens {
		account, err := w.CreateMasterAccount(w.ctx, password, token, w.walletDB.db)
		if err != nil {
			return fmt.Errorf("error creating %s account: %v", token, err)
		}

		w.Accounts[token] = account
	}

	return nil
}

func (w *Wallet) CreateMasterAccount(ctx context.Context, password, token string, db *sql.DB) (masterAccount, error) {
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing master public key: %v", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(ctx, password, pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error retrieving key from DB: %v", err)
	}

	factory, ok := masterAccountFactories[token]
	if !ok {
		return nil, fmt.Errorf("unsupported token type: %s", token)
	}

	masterAccount, err := factory(ctx, masterKey, db)
	if err != nil {
		return nil, fmt.Errorf("error creating %s account: %v", token, err)

	}

	return masterAccount, nil
}

func createETHAccount(ctx context.Context, masterKey *bip32.Key, db *sql.DB) (masterAccount, error) {

	ethAccount, err := eth.NewETHAccount(ctx, masterKey, "ETH", db)
	if err != nil {
		return nil, fmt.Errorf("error creating ETH account: %v", err)
	}

	return ethAccount, nil
}

func (w *Wallet) GetBalance(token string) (float64, error) {
	masterAccount, ok := w.Accounts[token]
	if !ok {
		return 0, fmt.Errorf("token not found: %s", token)
	}

	hexBalance, err := masterAccount.RetrieveBalance(0)
	if err != nil {
		return 0, fmt.Errorf("error retrieving balance: %v", err)
	}

	balance, err := eth.HexToEther(hexBalance)
	if err != nil {
		return 0, fmt.Errorf("error converting balance: %v", err)
	}

	return strconv.ParseFloat(balance, 64)
}

func (w *Wallet) EstimateGas(token, to, value string, accountIndex int) (string, error) {
	masterAccount, ok := w.Accounts[token]
	if !ok {
		return "", fmt.Errorf("token not found: %s", token)
	}
	gasPrice, err := masterAccount.EstimateGas(to, value, accountIndex)
	if err != nil {
		return "", fmt.Errorf("error estimating gad price for token %s : %v", token, err)
	}

	return gasPrice, nil
}

func (w *Wallet) SendTransaction(token, password, to, value string, accountIndex int) (bool, error) {
	masterAccount, ok := w.Accounts[token]
	if !ok {
		return false, fmt.Errorf("token not found: %s", token)
	}

	dbCtx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return false, fmt.Errorf("error serializing master public key: %v", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(dbCtx, password, pubKeyHex)
	if err != nil {
		return false, fmt.Errorf("error retrieving key from DB: %v", err)
	}

	from, err := masterAccount.GetAddress(accountIndex)
	if err != nil {
		return false, fmt.Errorf("error getting %s account address for index %d : %v", token, accountIndex, err)
	}

	transactionHash, err := masterAccount.SendTransaction(to, value, masterKey, accountIndex)
	if err != nil {
		return false, fmt.Errorf("failed to process %s transaction %v", token, err)
	}

	var status string
	if transactionHash == "0x0000000000000000000000000000000000000000000000000000000000000000" {
		status = "PENDING"
	} else {
		status = "COMPLETED"
	}

	now := time.Now().UTC()
	isoDate := now.Format(time.RFC3339)

	err = w.walletDB.SaveTransactionInDB(dbCtx, from, to, value, status, token, isoDate)
	if err != nil {
		return true, fmt.Errorf("error saving transaction into DB: %v", err)
	}

	return true, nil
}

func (w *Wallet) GetTransactions() ([]WalletTransaction, error) {
	dbCtx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()

	return w.walletDB.GetTransactions(dbCtx)
}
