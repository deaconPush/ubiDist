package hdwallet_test

import (
	"context"
	"encoding/hex"
	"testing"
	"time"
	"wallet/internal/hdwallet"
	"wallet/internal/utils"

	_ "modernc.org/sqlite"

	"github.com/stretchr/testify/require"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func TestWalletStorageOperations(t *testing.T) {
	ctx := context.Background()

	ws, err := hdwallet.NewWalletStorage(ctx, ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		ws.Close()
	})

	assertWalletExistence(ctx, t, ws, false)
	password := "password"
	wrongPassword := "wrong_password"
	pubKeyHex, encryptedMasterKeyHex := generateWallet(t, password)
	err = ws.SaveRootKeyToDB(ctx, pubKeyHex, encryptedMasterKeyHex)
	require.NoError(t, err)

	t.Run("Valid password retrieves correct root key", func(t *testing.T) {
		assertRootKeyRetrieval(ctx, t, ws, password, pubKeyHex, encryptedMasterKeyHex)
	})

	t.Run("Invalid password fails to retrieve root key", func(t *testing.T) {
		assertRootKeyRetrievalError(ctx, t, ws, wrongPassword, pubKeyHex)
	})

	t.Run("Test password validation", func(t *testing.T) {
		ok, err := ws.ValidatePassword(ctx, pubKeyHex, password)
		require.NoError(t, err)
		require.Equal(t, true, ok)

		ok, err = ws.ValidatePassword(ctx, pubKeyHex, wrongPassword)
		require.Error(t, err)
		require.Equal(t, false, ok)
	})
}

func generateWallet(t testing.TB, password string) (string, []byte) {
	t.Helper()
	entropy, err := bip39.NewEntropy(128)
	require.NoError(t, err)

	mnemonic, err := bip39.NewMnemonic(entropy)
	require.NoError(t, err)

	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	require.NoError(t, err)

	pubKeyData, err := masterKey.PublicKey().Serialize()
	require.NoError(t, err)
	pubKeyHex := hex.EncodeToString(pubKeyData)

	masterKeyData, err := masterKey.Serialize()
	require.NoError(t, err)
	masterKeyHex := hex.EncodeToString(masterKeyData)

	encryptedMasterKeyHex, err := utils.Encrypt([]byte(password), []byte(masterKeyHex))
	require.NoError(t, err)

	return pubKeyHex, encryptedMasterKeyHex
}

func assertWalletExistence(ctx context.Context, t testing.TB, storage *hdwallet.WalletStorage, want bool) {
	t.Helper()

	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	got, err := storage.WalletExists(dbCtx)
	require.NoError(t, err)
	require.Equal(t, want, got, "Wallet existence mismatch")
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
	require.NoError(t, err)
	require.NotNil(t, retrievedKey)

	masterKey, err := decryptMasterKey(password, encryptedMasterKeyHex)
	require.NoError(t, err)

	require.Equal(t, masterKey.String(), retrievedKey.String(), "Retrieved key does not match original")
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
	require.Error(t, err)
	require.Nil(t, retrievedKey)
}
