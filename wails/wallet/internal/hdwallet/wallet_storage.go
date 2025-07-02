package hdwallet

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	_ "modernc.org/sqlite"
)

type WalletStorage struct {
	db *sql.DB
}

type WalletTransaction struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Status    string `json:"status"`
	Value     string `json:"value"`
	Token     string `json:"token"`
	CreatedAt string `json:"createdAt"`
}

func NewWalletStorage(ctx context.Context, filePath string) (*WalletStorage, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS wallets (publicKey TEXT PRIMARY KEY, masterKey TEXT)")
	if err != nil {
		return nil, fmt.Errorf("error creating wallets table: %w", err)
	}

	_, err = db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS transactions (
		sender TEXT,
		recipient TEXT,
		value TEXT,
		status TEXT,
		token TEXT,
		createdAt TEXT
	)`)
	if err != nil {
		return nil, fmt.Errorf("error creating wallets table: %w", err)
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

func (ws *WalletStorage) GetTransactions(ctx context.Context) ([]WalletTransaction, error) {
	var transactions []WalletTransaction
	rows, err := ws.db.QueryContext(ctx, "SELECT * FROM transactions ORDER BY createdAt DESC")
	if err != nil {
		return nil, fmt.Errorf("error retrieving transactions from db: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var transaction WalletTransaction
		if err := rows.Scan(&transaction.Sender, &transaction.Recipient, &transaction.Value, &transaction.Status, &transaction.Token, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("error parsing db transaction data: %w", err)
		}

		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error retrieving transaction rows from db: %w", err)
	}

	return transactions, nil
}

func (ws *WalletStorage) ValidatePassword(ctx context.Context, pubKeyHex, password string) (bool, error) {
	var encryptedKeyData string
	err := ws.db.QueryRowContext(ctx, "SELECT masterKey FROM wallets WHERE publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no rows returned")
		}
		return false, fmt.Errorf("error querying database: %w", err)
	}

	_, err = utils.Decrypt([]byte(password), []byte(encryptedKeyData))
	if err != nil {
		return false, fmt.Errorf("password is invalid: %w", err)
	}

	return true, nil
}

func (ws *WalletStorage) SaveRootKeyToDB(ctx context.Context, pubKeyHex string, encryptedMasterKey []byte) error {
	result, err := ws.db.ExecContext(
		ctx,
		`INSERT INTO wallets (publicKey, masterKey) VALUES (?, ?)`,
		pubKeyHex,
		encryptedMasterKey,
	)
	if err != nil {
		return fmt.Errorf("error saving HDKey: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error inserting record into DB: %w", err)
	}

	return nil
}

func (ws *WalletStorage) RetrieveRootKeyFromDB(ctx context.Context, password, pubKeyHex string) (*bip32.Key, error) {
	var encryptedKeyData string
	err := ws.db.QueryRowContext(ctx, "SELECT masterKey FROM wallets WHERE publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no rows returned")
		}
		return nil, fmt.Errorf("error querying database: %w", err)
	}

	keyDataHex, err := utils.Decrypt([]byte(password), []byte(encryptedKeyData))
	if err != nil {
		return nil, fmt.Errorf("error decrypting master key %w", err)
	}

	keyData, err := hex.DecodeString(string(keyDataHex))
	if err != nil {
		return nil, fmt.Errorf("error decoding key data: %w", err)
	}

	masterKey, err := bip32.Deserialize(keyData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing HDKey from wallet file: %w", err)
	}

	return masterKey, nil
}

func (ws *WalletStorage) RetrievePublicKeyFromDB(ctx context.Context) (*bip32.Key, error) {
	var pubKeyHex string
	err := ws.db.QueryRowContext(ctx, "SELECT publicKey FROM wallets").Scan(&pubKeyHex)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no rows returned")
		}
		return nil, fmt.Errorf("error querying database: %w", err)
	}

	pubKeyData, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return nil, fmt.Errorf("error decoding public key data: %w", err)
	}

	pubKey, err := bip32.Deserialize(pubKeyData)
	if err != nil {
		return nil, fmt.Errorf("error deserializing public key: %w", err)
	}

	return pubKey, nil
}

func (ws *WalletStorage) SaveTransactionInDB(ctx context.Context, from, to, value, status, token, date string) error {
	result, err := ws.db.ExecContext(
		ctx,
		`INSERT INTO transactions 
		(sender, recipient, value, status, token, createdAt) 
	 VALUES (?, ?, ?, ?, ?, ?)`,
		from,
		to,
		value,
		status,
		token,
		date,
	)
	if err != nil {
		return fmt.Errorf("error saving transaction into DB: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error retrieving rows affected: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("error inserting record into DB: %w", err)
	}

	return nil
}

func (ws *WalletStorage) Close() error {
	if ws.db != nil {
		return ws.db.Close()
	}
	return nil
}
