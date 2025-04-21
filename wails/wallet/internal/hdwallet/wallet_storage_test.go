package hdwallet

import (
	"context"
	"encoding/hex"
	"reflect"
	"testing"
	"time"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func TestWalletStorageOperations(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	ws, err := NewWalletStorage(":memory:", ctx)
	if err != nil {
		t.Fatalf("Failed to create database service: %v", err)
	}

	defer cancel()
	assertWalletExistence(t, ws, false, ctx)

	password := "password"
	pubKeyHex, encryptedMasterKeyHex, err := generateWallet(t, password)
	if err != nil {
		t.Errorf("Failed to generate wallet: %v", err)
	}

	err = ws.SaveRootKeyToDB(ctx, password, pubKeyHex, encryptedMasterKeyHex)
	if err != nil {
		t.Errorf("Failed to save root key to DB: %v", err)
	}

	cases := []struct {
		name   string
		testFn func(t *testing.T, storage *WalletStorage)
	}{
		{
			name: "Valid password retrieves correct root key",
			testFn: func(t *testing.T, storage *WalletStorage) {
				assertRootKeyRetrieval(t, storage, password, pubKeyHex, encryptedMasterKeyHex, ctx)
			},
		},
		{
			name: "Invalid password fails to retrieve root key",
			testFn: func(t *testing.T, storage *WalletStorage) {
				assertRootKeyRetrievalError(t, storage, "wrong_password", pubKeyHex, ctx)
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

func generateWallet(t testing.TB, password string) (string, []byte, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		t.Errorf("Failed to generate entropy: %v", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		t.Errorf("Failed to generate mnemonic: %v", err)
	}
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		t.Errorf("Failed to create master key: %v", err)
	}
	pubKeyData, err := masterKey.PublicKey().Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize public key: %v", err)
	}
	pubKeyHex := hex.EncodeToString(pubKeyData)
	masterKeyData, err := masterKey.Serialize()
	if err != nil {
		t.Errorf("Failed to serialize master key: %v", err)
	}
	masterKeyHex := hex.EncodeToString(masterKeyData)
	encryptedMasterKeyHex, err := utils.Encrypt([]byte(password), []byte(masterKeyHex))
	if err != nil {
		t.Errorf("Failed to encrypt master key: %v", err)
	}
	return pubKeyHex, encryptedMasterKeyHex, nil
}

func assertWalletExistence(t testing.TB, storage *WalletStorage, want bool, ctx context.Context) {
	t.Helper()
	got, err := storage.WalletExists(ctx)
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

func assertRootKeyRetrieval(t testing.TB, storage *WalletStorage, password, pubKeyHex string, encryptedMasterKeyHex []byte, ctx context.Context) {
	t.Helper()
	retrievedKey, err := storage.RetrieveRootKeyFromDB(ctx, password, pubKeyHex)
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

func assertRootKeyRetrievalError(t testing.TB, storage *WalletStorage, password, pubKeyHex string, ctx context.Context) {
	t.Helper()
	retrievedKey, err := storage.RetrieveRootKeyFromDB(ctx, password, pubKeyHex)
	if err == nil {
		t.Fatalf("Expected an error, but got a valid key: %v", retrievedKey)
	}
}
