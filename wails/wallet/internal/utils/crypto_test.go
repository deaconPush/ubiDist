package utils_test

import (
	"encoding/hex"
	"testing"
	"wallet/internal/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func TestEncryption(t *testing.T) {
	t.Run("Encrypt and decrypt with valid password and data", func(t *testing.T) {
		password := "password"
		plain := []byte("data")
		ciphertext, err := utils.Encrypt([]byte(password), plain)
		require.NoError(t, err)

		decrypted, err := utils.Decrypt([]byte(password), ciphertext)
		require.NoError(t, err)
		require.Equal(t, plain, decrypted)
	})

	t.Run("Decrypt with invalid password should fail", func(t *testing.T) {
		password := "password"
		invalidPassword := "invalid"
		plain := []byte("data")
		ciphertext, err := utils.Encrypt([]byte(password), plain)
		require.NoError(t, err)

		_, err = utils.Decrypt([]byte(invalidPassword), ciphertext)
		require.Error(t, err)
	})

	t.Run("Encrypt and decrypt a bip32 key", func(t *testing.T) {
		password := "password"
		seed, err := bip32.NewSeed()
		require.NoError(t, err)

		key, err := bip32.NewMasterKey(seed)
		require.NoError(t, err)

		serializedKey, err := key.Serialize()
		require.NoError(t, err)

		plain := []byte(hex.EncodeToString(serializedKey))
		ciphertext, err := utils.Encrypt([]byte(password), plain)
		require.NoError(t, err)

		decrypted, err := utils.Decrypt([]byte(password), ciphertext)
		require.NoError(t, err)
		require.Equal(t, plain, decrypted)
	})
}

func TestAddressValidation(t *testing.T) {
	cases := []struct {
		name    string
		address string
		token   string
		isValid bool
	}{
		{
			name:    "All lowercase ETH address",
			address: "0xde709f2102306220921060314715629080e2fb77",
			token:   "ETH",
			isValid: true,
		},
		{
			name:    "All uppercase ETH address",
			address: "0XDE709F2102306220921060314715629080E2FB77",
			token:   "ETH",
			isValid: true,
		},
		{
			name:    "Invalid checksum",
			address: "0x52908400098527886e0F7030069857D2E4169EE7",
			token:   "ETH",
			isValid: false,
		},
		{
			name:    "Valid checksum",
			address: "0x27b1FDB04752BBC536007A920D24ACB045561c26",
			token:   "ETH",
			isValid: true,
		},
		{
			name:    "Invalid address",
			address: "0x1234567890abcdef1234567890abcdef1234567890",
			token:   "ETH",
			isValid: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := utils.ValidateAddress(tc.address, tc.token)
			require.Equal(t, tc.isValid, result)
		})
	}
}

func TestAccountKeyDerivation(t *testing.T) {
	mnemonic := "test test test test test test test test test test test junk"
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	require.NoError(t, err)

	cases := []struct {
		name            string
		coinType        string
		accountIndex    int
		expectedAddress string
	}{
		{
			name:            "1st ETH account",
			coinType:        "ETH",
			accountIndex:    0,
			expectedAddress: "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		},
		{
			name:            "4th ETH account",
			coinType:        "ETH",
			accountIndex:    4,
			expectedAddress: "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
		},
		{
			name:            "8th ETH account",
			coinType:        "ETH",
			accountIndex:    8,
			expectedAddress: "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			childKey, err := utils.DeriveKeyForAccount(masterKey, tc.coinType, tc.accountIndex)
			require.NoError(t, err)

			privateKey, err := crypto.ToECDSA(childKey.Key)
			require.NoError(t, err)

			address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
			require.Equal(t, tc.expectedAddress, address)
		})
	}
}
