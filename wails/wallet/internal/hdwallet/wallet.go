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
	GetAllAccounts() (map[int]string, error)
}

type masterAccountFactory func(ctx context.Context, masterKey *bip32.Key, db *sql.DB) (masterAccount, error)

var masterAccountFactories = map[string]masterAccountFactory{
	"ETH": createETHAccount,
}

func CreateWallet(ctx context.Context, password string, ws *WalletStorage) (*Wallet, string, error) {
	mnemonic, err := utils.GenerateMnemonic()
	if err != nil {
		return nil, "", fmt.Errorf("error generating mnemonic: %w", err)
	}

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, "", fmt.Errorf("error recovering master key from seed: %w", err)
	}

	err = storeMasterKey(ctx, ws, password, masterKey)
	if err != nil {
		return nil, "", fmt.Errorf("error storing master key into local db: %w", err)
	}

	return &Wallet{
		publicKey: masterKey.PublicKey(),
		Accounts:  make(map[string]masterAccount),
		walletDB:  ws,
		ctx:       ctx,
	}, mnemonic, nil
}

func RestoreWallet(ctx context.Context, password string, mnemonic string, ws *WalletStorage) (*Wallet, error) {
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("error recovering master key from seed: %w", err)
	}

	err = storeMasterKey(ctx, ws, password, masterKey)
	if err != nil {
		return nil, fmt.Errorf("error storing master key: %w", err)
	}

	return &Wallet{publicKey: masterKey.PublicKey(), Accounts: make(map[string]masterAccount), walletDB: ws, ctx: ctx}, nil
}

func RecoverWallet(ctx context.Context, password string, ws *WalletStorage) (*Wallet, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	pubKey, err := ws.RetrievePublicKeyFromDB(dbCtx)
	defer cancel()
	if err != nil {
		return nil, fmt.Errorf("error retrieving public key from DB: %w", err)
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
		return fmt.Errorf("error serializing master Key: %w", err)
	}

	pubKeyData, err := masterKey.PublicKey().Serialize()
	if err != nil {
		return fmt.Errorf("error serializing master public key: %w", err)
	}

	masterKeyHex := hex.EncodeToString(masterKeyData)
	encryptedMasterKey, err := utils.Encrypt([]byte(password), []byte(masterKeyHex))
	if err != nil {
		return fmt.Errorf("error encrypting data: %w", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	err = ws.SaveRootKeyToDB(dbCtx, pubKeyHex, encryptedMasterKey)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %w", err)
	}

	return nil
}

func (w *Wallet) Initialize(tokens []string, password string) error {
	for _, token := range tokens {
		account, err := w.CreateMasterAccount(w.ctx, password, token, w.walletDB.db)
		if err != nil {
			return fmt.Errorf("error creating %s account: %w", token, err)
		}

		w.Accounts[token] = account
	}

	return nil
}

func (w *Wallet) CreateMasterAccount(ctx context.Context, password, token string, db *sql.DB) (masterAccount, error) {
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return nil, fmt.Errorf("error serializing master public key: %w", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(ctx, password, pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error retrieving key from DB: %w", err)
	}

	factory, ok := masterAccountFactories[token]
	if !ok {
		return nil, fmt.Errorf("unsupported token type: %s", token)
	}

	masterAcct, err := factory(ctx, masterKey, db)
	if err != nil {
		return nil, fmt.Errorf("error creating %s account: %w", token, err)
	}

	return masterAcct, nil
}

func createETHAccount(ctx context.Context, masterKey *bip32.Key, db *sql.DB) (masterAccount, error) {
	ethAccount, err := eth.NewETHAccount(ctx, masterKey, "ETH", db)
	if err != nil {
		return nil, fmt.Errorf("error creating ETH account: %w", err)
	}

	return ethAccount, nil
}

func (w *Wallet) GetAccountAddress(token string, accountIndex int) (string, error) {
	masterAcc, ok := w.Accounts[token]
	if !ok {
		return "", fmt.Errorf("token not found: %s", token)
	}

	address, err := masterAcc.GetAddress(accountIndex)
	if err != nil {
		return "", fmt.Errorf("error getting %s account address for index %d : %w", token, accountIndex, err)
	}

	return address, nil
}

func (w *Wallet) GetAllAccounts(token string) (map[int]string, error) {
	masterAcc, ok := w.Accounts[token]
	if !ok {
		return nil, fmt.Errorf("token not found: %s", token)
	}

	accounts, err := masterAcc.GetAllAccounts()
	if err != nil {
		return nil, fmt.Errorf("error retrieving accounts for token: %s : %w", token, err)
	}

	return accounts, nil
}

func (w *Wallet) GetBalance(token string, accountIndex int) (float64, error) {
	masterAcc, ok := w.Accounts[token]
	if !ok {
		return 0, fmt.Errorf("token not found: %s", token)
	}

	hexBalance, err := masterAcc.RetrieveBalance(accountIndex)
	if err != nil {
		return 0, fmt.Errorf("error retrieving balance: %w", err)
	}

	balance, err := eth.HexToEther(hexBalance)
	if err != nil {
		return 0, fmt.Errorf("error converting balance: %w", err)
	}

	return strconv.ParseFloat(balance, 64)
}

func (w *Wallet) EstimateGas(token, to, value string, accountIndex int) (string, error) {
	masterAcc, ok := w.Accounts[token]
	if !ok {
		return "", fmt.Errorf("token not found: %s", token)
	}
	gasPrice, err := masterAcc.EstimateGas(to, value, accountIndex)
	if err != nil {
		return "", fmt.Errorf("error estimating gad price for token %s : %w", token, err)
	}

	return gasPrice, nil
}

func (w *Wallet) SendTransaction(token, password, to, value string, accountIndex int) (bool, error) {
	masterAcc, ok := w.Accounts[token]
	if !ok {
		return false, fmt.Errorf("token not found: %s", token)
	}

	dbCtx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()
	pubKeyData, err := w.publicKey.Serialize()
	if err != nil {
		return false, fmt.Errorf("error serializing master public key: %w", err)
	}

	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKey, err := w.walletDB.RetrieveRootKeyFromDB(dbCtx, password, pubKeyHex)
	if err != nil {
		return false, fmt.Errorf("error retrieving key from DB: %w", err)
	}

	from, err := masterAcc.GetAddress(accountIndex)
	if err != nil {
		return false, fmt.Errorf("error getting %s account address for index %d : %w", token, accountIndex, err)
	}

	transactionHash, err := masterAcc.SendTransaction(to, value, masterKey, accountIndex)
	if err != nil {
		return false, fmt.Errorf("failed to process %s transaction %w", token, err)
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
		return true, fmt.Errorf("error saving transaction into DB: %w", err)
	}

	return true, nil
}

func (w *Wallet) GetTransactions() ([]WalletTransaction, error) {
	dbCtx, cancel := context.WithTimeout(w.ctx, 5*time.Second)
	defer cancel()

	return w.walletDB.GetTransactions(dbCtx)
}
