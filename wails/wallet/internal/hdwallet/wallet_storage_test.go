package hdwallet

import (
	"context"
	"encoding/hex"
	"reflect"
	"testing"
	"wallet/internal/utils"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func TestWalletStorageOperations(t *testing.T) {
	dbService, err := utils.NewDatabaseService(context.Background(), ":memory:")
	if err != nil {
		t.Fatalf("Failed to create database service: %v", err)
	}
	defer dbService.GetDB().Close()

	storage := NewWalletStorage(dbService.GetDB())
	assertWalletExistence(t, storage, context.Background(), false)

	password := "password"
	pubKeyHex, encryptedMasterKeyHex, err := generateWallet(t, password)
	if err != nil {
		t.Errorf("Failed to generate wallet: %v", err)
	}

	err = storage.SaveRootKeyToDB(password, pubKeyHex, encryptedMasterKeyHex)
	if err != nil {
		t.Errorf("Failed to save root key to DB: %v", err)
	}
	assertWalletExistence(t, storage, context.Background(), true)
	assertRootKeyRetrieval(t, storage, password, pubKeyHex, encryptedMasterKeyHex)
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

func assertWalletExistence(t testing.TB, storage *WalletStorage, ctx context.Context, want bool) {
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

func assertRootKeyRetrieval(t testing.TB, storage *WalletStorage, password, pubKeyHex string, encryptedMasterKeyHex []byte) {
	t.Helper()
	retrievedKey, err := storage.RetrieveRootKeyFromDB(password, pubKeyHex)
	if err != nil {
		t.Fatalf("Failed to retrieve root key: %v", err)
	}

	if retrievedKey == nil {
		t.Fatalf("Expected a valid key, but got nil")
	}

	// Decrypt and compare with the original master key
	masterKey, err := decryptMasterKey(password, encryptedMasterKeyHex)
	if err != nil {
		t.Fatalf("Failed to decrypt master key: %v", err)
	}

	if !reflect.DeepEqual(retrievedKey, masterKey) {
		t.Errorf("Retrieved key does not match the original")
	}
}
