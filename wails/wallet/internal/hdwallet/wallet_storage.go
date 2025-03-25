package hdwallet

import (
	"context"
	"database/sql"
	"fmt"
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

func (ws *WalletStorage) SaveRootKeyToDB(password string, pubKeyHex string, encryptedMasterKey []byte) error {
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

func (ws *WalletStorage) retrieveRootKeyFromDB() {}
