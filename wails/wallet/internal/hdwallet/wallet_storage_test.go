package hdwallet_test

import (
	"context"
	"encoding/hex"
	"reflect"
	"testing"
	"time"
	"wallet/internal/hdwallet"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func TestWalletStorageOperations(t *testing.T) {
	ctx := context.Background()
	ws, err := hdwallet.NewWalletStorage(ctx, ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database service: %v", err)
	}

	assertWalletExistence(ctx, t, ws, false)

	password := "password"
	pubKeyHex, encryptedMasterKeyHex := generateWallet(t, password)
	err = ws.SaveRootKeyToDB(ctx, pubKeyHex, encryptedMasterKeyHex)
	if err != nil {
		t.Fatalf("Failed to save root key to DB: %v", err)
	}

	cases := []struct {
		name   string
		testFn func(t *testing.T, storage *hdwallet.WalletStorage)
	}{
		{
			name: "Valid password retrieves correct root key",
			testFn: func(t *testing.T, storage *hdwallet.WalletStorage) {
				assertRootKeyRetrieval(ctx, t, storage, password, pubKeyHex, encryptedMasterKeyHex)
			},
		},
		{
			name: "Invalid password fails to retrieve root key",
			testFn: func(t *testing.T, storage *hdwallet.WalletStorage) {
				assertRootKeyRetrievalError(ctx, t, storage, "wrong_password", pubKeyHex)
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.testFn(t, ws)
		})
	}
	defer ws.Close()
}

func generateWallet(t testing.TB, password string) (string, []byte) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Fatalf("Failed to generate entropy: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Fatalf("Failed to generate mnemonic: %v", err)
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		t.Fatalf("Failed to create master key: %v", err)
	}
	pubKeyData, err := masterKey.PublicKey().Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize public key: %v", err)
	}
	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKeyData, err := masterKey.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize master key: %v", err)
	}
	masterKeyHex := hex.EncodeToString(masterKeyData)
	encryptedMasterKeyHex, err := utils.Encrypt([]byte(password), []byte(masterKeyHex))
	if err != nil {
		t.Fatalf("Failed to encrypt master key: %v", err)
	}
	return pubKeyHex, encryptedMasterKeyHex
}

func assertWalletExistence(ctx context.Context, t testing.TB, storage *hdwallet.WalletStorage, want bool) {
	t.Helper()
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	got, err := storage.WalletExists(dbCtx)
	if err != nil {
		t.Fatalf("Error checking wallet existence: %v", err)
	}
	if got != want {
		t.Errorf("Wallet existence mismatch: expected %v, got %v", want, got)
	}
}

func decryptMasterKey(password string, encryptedMasterKeyHex []byte) (*bip32.Key, error) {
	masterKeyHex, err := utils.Decrypt([]byte(password), encryptedMasterKeyHex)
	if err != nil {
		return nil, err
	}

	masterKeyData, err := hex.DecodeString(string(masterKeyHex))
	if err != nil {
		return nil, err
	}

	return bip32.Deserialize(masterKeyData)
}

func assertRootKeyRetrieval(
	ctx context.Context,
	t testing.TB,
	storage *hdwallet.WalletStorage,
	password, pubKeyHex string,
	encryptedMasterKeyHex []byte,
) {
	t.Helper()
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	retrievedKey, err := storage.RetrieveRootKeyFromDB(dbCtx, password, pubKeyHex)
	if err != nil {
		t.Errorf("Failed to retrieve root key: %v", err)
	}

	if retrievedKey == nil {
		t.Errorf("Expected a valid key, but got nil")
	}

	// Decrypt and compare with the original master key
	masterKey, err := decryptMasterKey(password, encryptedMasterKeyHex)
	if err != nil {
		t.Errorf("Failed to decrypt master key: %v", err)
	}

	if !reflect.DeepEqual(retrievedKey, masterKey) {
		t.Errorf("Retrieved key does not match the original")
	}
}

func assertRootKeyRetrievalError(
	ctx context.Context,
	t testing.TB,
	storage *hdwallet.WalletStorage,
	password,
	pubKeyHex string) {
	t.Helper()
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	retrievedKey, err := storage.RetrieveRootKeyFromDB(dbCtx, password, pubKeyHex)
	if err == nil {
		t.Fatalf("Expected an error, but got a valid key: %v", retrievedKey)
	}
}
