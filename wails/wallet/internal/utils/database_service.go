package utils

import (
	"context"
	"database/sql"
	"fmt"
)

type DatabaseService struct {
	db *sql.DB
}

func NewDatabaseService(ctx context.Context) (*DatabaseService, error) {
	db, err := sql.Open("sqlite3", "wallet.db")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	service := &DatabaseService{db: db}
	if err := service.initDatabase(ctx); err != nil {
		return nil, fmt.Errorf("error initializing database: %v", err)
	}

	return service, nil
}

func (d *DatabaseService) initDatabase(ctx context.Context) error {
	_, err := d.db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS wallets (publicKey TEXT PRIMARY KEY, masterKey TEXT)")
	return err
}

func (d *DatabaseService) GetDB() *sql.DB {
	return d.db
}
