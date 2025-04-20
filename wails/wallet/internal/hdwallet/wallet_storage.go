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
	db  *sql.DB
	ctx context.Context
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

	return &WalletStorage{db: db, ctx: ctx}, nil
}

func (ws *WalletStorage) WalletExists() (bool, error) {
	var count int
	err := ws.db.QueryRowContext(ws.ctx, "SELECT COUNT(*) FROM wallets").Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ws *WalletStorage) SaveRootKeyToDB(password, pubKeyHex string, encryptedMasterKey []byte) error {
	result, err := ws.db.ExecContext(ws.ctx, "INSERT INTO wallets (publicKey, masterKey) VALUES (?, ?)", pubKeyHex, encryptedMasterKey)
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

func (ws *WalletStorage) RetrieveRootKeyFromDB(password, pubKeyHex string) (*bip32.Key, error) {
	var encryptedKeyData string
	err := ws.db.QueryRowContext(ws.ctx, "SELECT masterKey FROM wallets WHERE publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
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

func (ws *WalletStorage) RetrievePublicKeyFromDB() (*bip32.Key, error) {
	var pubKeyHex string
	err := ws.db.QueryRowContext(ws.ctx, "SELECT publicKey FROM wallets").Scan(&pubKeyHex)
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

func (w *WalletStorage) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}
