package hdwallet

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	_ "modernc.org/sqlite"
)

type WalletStorage struct {
	db *sql.DB
}

type WalletTransaction struct {
	Sender    string
	Recipient string
	Status    string
}

func NewWalletStorage(filePath string, ctx context.Context) (*WalletStorage, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS wallets (publicKey TEXT PRIMARY KEY, masterKey TEXT)")
	if err != nil {
		return nil, fmt.Errorf("error creating wallets table: %v", err)
	}

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS transactions (sender TEXT PRIMARY KEY, recipient TEXT, status TEXT)")
	if err != nil {
		return nil, fmt.Errorf("error creating wallets table: %v", err)
	}

	return &WalletStorage{db: db}, nil
}

func (ws *WalletStorage) WalletExists(ctx context.Context) (bool, error) {
	var count int
	err := ws.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wallets").Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ws *WalletStorage) GetTransactions(ctx context.Context, pubKey string) ([]WalletTransaction, error) {
	var transactions []WalletTransaction
	rows, err := ws.db.QueryContext(ctx, "SELECT * FROM transactions where sender=?", pubKey)
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions from db: %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var transaction WalletTransaction
		if err := rows.Scan(&transaction.Sender, &transaction.Recipient, &transaction.Status); err != nil {
			return nil, fmt.Errorf("error parsing db transaction data")
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error retrieving transaction rows from db: %v", err)
	}

	return transactions, nil
}

func (ws *WalletStorage) ValidatePassword(ctx context.Context, pubKeyHex, password string) (bool, error) {
	var encryptedKeyData string
	err := ws.db.QueryRowContext(ctx, "SELECT masterKey FROM wallets WHERE publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("no rows returned")
		}
		return false, fmt.Errorf("error querying database: %v", err)
	}

	_, err = utils.Decrypt([]byte(password), []byte(encryptedKeyData))
	if err != nil {
		return false, fmt.Errorf("password is invalid: %v", err)
	}

	return true, nil
}

func (ws *WalletStorage) SaveRootKeyToDB(ctx context.Context, password, pubKeyHex string, encryptedMasterKey []byte) error {
	result, err := ws.db.ExecContext(ctx, "INSERT INTO wallets (publicKey, masterKey) VALUES (?, ?)", pubKeyHex, encryptedMasterKey)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %v", err)
	}

	if rows != 1 {
		return fmt.Errorf("error inserting record into DB: %v", err)
	}

	return nil
}

func (ws *WalletStorage) RetrieveRootKeyFromDB(ctx context.Context, password, pubKeyHex string) (*bip32.Key, error) {
	var encryptedKeyData string
	err := ws.db.QueryRowContext(ctx, "SELECT masterKey FROM wallets WHERE publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows returned")
		}
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

func (ws *WalletStorage) RetrievePublicKeyFromDB(ctx context.Context) (*bip32.Key, error) {
	var pubKeyHex string
	err := ws.db.QueryRowContext(ctx, "SELECT publicKey FROM wallets").Scan(&pubKeyHex)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no rows returned")
		}
		return nil, fmt.Errorf("error querying database: %v", err)
	}

	pubKeyData, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding public key data: %v", err)
	}

	pubKey, err := bip32.Deserialize(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing public key: %v", err)
	}

	return pubKey, nil
}

func (ws *WalletStorage) SaveTransactionInDB(ctx context.Context, from, to, value, status string) error {
	result, err := ws.db.ExecContext(ctx, "INSERT INTO transactions (from, to, value, status) VALUES(?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error saving transaction into DB: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %v", err)
	}

	if rows != 1 {
		return fmt.Errorf("error inserting record into DB: %v", err)
	}

	return nil
}

func (w *WalletStorage) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}
