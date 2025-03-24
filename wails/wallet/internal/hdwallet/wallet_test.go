package hdwallet

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func TestKeyPersistence(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := initDB(t, ctx)
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}

}
func initDB(t testing.TB, ctx context.Context) (*sql.DB, error) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}

	db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS wallets (publicKey TEXT PRIMARY KEY, masterKey TEXT)")
	return db, nil
}
