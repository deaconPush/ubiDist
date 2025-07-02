package eth

import (
	"context"
	"database/sql"
	"fmt"
)

type AccountStorage struct {
	db *sql.DB
}

func NewAccountStorage(ctx context.Context, db *sql.DB) (*AccountStorage, error) {
	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS ethAccounts (address TEXT, accountIndex INTEGER PRIMARY KEY)")
	if err != nil {
		return nil, fmt.Errorf("error creating wallets table: %w", err)
	}

	return &AccountStorage{db: db}, nil
}

func (a *AccountStorage) AccountsExist(ctx context.Context) (bool, error) {
	var count int
	err := a.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM ethAccounts").Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (a *AccountStorage) SaveAccounts(ctx context.Context, accounts []string) error {
	stmt, err := a.db.PrepareContext(ctx, "INSERT INTO ethAccounts (address, accountIndex) VALUES (?, ?)")
	if err != nil {
		return fmt.Errorf("error preparing statement for inserting eth accounts: %w", err)
	}

	defer stmt.Close()

	for i := 0; i < len(accounts); i++ {
		_, err = stmt.ExecContext(ctx, accounts[i], i)
		if err != nil {
			return fmt.Errorf("error inserting eth account %d : %w", i, err)
		}
	}

	return nil
}

func (a *AccountStorage) GetAccountAddress(ctx context.Context, accountIndex int) (string, error) {
	var address string
	err := a.db.QueryRowContext(ctx, "SELECT address FROM ethAccounts where accountIndex =?", accountIndex).Scan(&address)

	if err != nil {
		return "", fmt.Errorf("error retrieving ETH account %d from DB: %w", accountIndex, err)
	}

	return address, nil
}

func (a *AccountStorage) GetAllAccounts(ctx context.Context) (map[int]string, error) {
	rows, err := a.db.QueryContext(ctx, "SELECT address, accountIndex FROM ethAccounts LIMIT 10")
	accounts := make(map[int]string)
	if err != nil {
		return nil, fmt.Errorf("error querying ethAccounts: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var accountIndex int
		var address string
		if err := rows.Scan(&address, &accountIndex); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		accounts[accountIndex] = address
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error retrieving ethAccounts rows from db: %w", err)
	}

	return accounts, nil
}
