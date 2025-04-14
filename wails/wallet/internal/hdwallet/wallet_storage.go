package hdwallet

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
)

type WalletStorage struct {
	db *sql.DB
}

func NewWalletStorage(db *sql.DB) *WalletStorage {
	return &WalletStorage{db: db}
}

func (ws *WalletStorage) WalletExists(ctx context.Context) (bool, error) {
	var count int
	err := ws.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM wallets").Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (ws *WalletStorage) SaveRootKeyToDB(password, pubKeyHex string, encryptedMasterKey []byte) error {
	result, err := ws.db.ExecContext(context.Background(), "INSERT INTO wallets (publicKey, masterKey) VALUES (?, ?)", pubKeyHex, encryptedMasterKey)
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
	err := ws.db.QueryRowContext(context.Background(), "SELECT masterKey FROM wallets WHERE publicKey=?", pubKeyHex).Scan(&encryptedKeyData)
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

func (ws *WalletStorage) RetrieveKeysFromDB(password string) (string, string, error) {
	var pubKeyHex string
	var encryptedRootHex string
	err := ws.db.QueryRowContext(context.Background(), "SELECT publicKey, masterKey FROM wallets").Scan(&pubKeyHex, &encryptedRootHex)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no rows returned")
		}
		return "", "", fmt.Errorf("error querying database: %v", err)
	}

	if err != nil {
		return "", "", fmt.Errorf("error decoding public key: %v", err)
	}

	return pubKeyHex, encryptedRootHex, nil
}
