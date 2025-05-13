package utils

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

func TestEncryption(t *testing.T) {
	t.Run("Encryption with valid password and data", func(t *testing.T) {
		password := "password"
		want := []byte("data")
		ciphertext, err := Encrypt([]byte(password), want)
		if err != nil {
			t.Errorf("Encryption failed: %v", err)
		}

		got, err := Decrypt([]byte(password), ciphertext)
		if err != nil {
			t.Errorf("Decryption failed: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Decrypted data does not match original data")
		}
	})

	t.Run("Encryption with invalid password", func(t *testing.T) {
		password := "password"
		want := []byte("data")
		ciphertext, err := Encrypt([]byte(password), want)
		if err != nil {
			t.Errorf("Encryption failed: %v", err)
		}

		invalidPassword := "invalid"
		_, err = Decrypt([]byte(invalidPassword), ciphertext)
		if err == nil {
			t.Errorf("Decryption should have failed")
		}
	})

	t.Run("Encryption of a bip32 key", func(t *testing.T) {
		password := "password"
		seed, err := bip32.NewSeed()
		if err != nil {
			t.Errorf("Error generating seed: %v", err)
		}

		key, err := bip32.NewMasterKey(seed)
		if err != nil {
			t.Errorf("Error generating master key: %v", err)
		}

		serializedKey, err := key.Serialize()
		if err != nil {
			t.Errorf("Error serializing key: %v", err)
		}

		want := []byte(hex.EncodeToString(serializedKey))
		ciphertext, err := Encrypt([]byte(password), want)
		if err != nil {
			t.Errorf("Encryption failed: %v", err)
		}

		got, err := Decrypt([]byte(password), ciphertext)
		if err != nil {
			t.Errorf("Decryption failed: %v", err)
		}

		assertCorrectValue(t, got, want)
	})
}

func TestChildKeyGeneration(t *testing.T) {
	mnemonic := "test test test test test test test test test test test junk"
	seed := bip39.NewSeed(mnemonic, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		t.Errorf("Error generating master key: %v", err)
	}

	cases := []struct {
		derivationPath string
		masterKey      string
		name           string
	}{
		{
			"m/44'/60'/0'/0/0",
			"ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80",
			"Generate ETH first key from default mnemonic",
		},
		{
			"m/44'/60'/0'/0/1",
			"59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d",
			"Generate ETH second key from default mnemonic",
		},
		{
			"m/44'/60'/0'/0/2",
			"5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a",
			"Generate ETH third key from default mnemonic",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			derivatedKey, err := DeriveChildKey(masterKey, tc.derivationPath)
			if err != nil {
				t.Errorf("Error deriving child key: %v", err)
			}

			got, err := crypto.ToECDSA(derivatedKey.Key)
			if err != nil {
				t.Errorf("Error converting key to ECDSA: %v", err)
			}

			want, err := crypto.HexToECDSA(tc.masterKey)
			if err != nil {
				t.Errorf("Error converting key to ECDSA: %v", err)
			}

			assertCorrectValue(t, got, want)
		})
	}
}

func TestAddressValidation(t *testing.T) {
	cases := []struct {
		address string
		name    string
		token   string
		isValid bool
	}{
		{
			address: "0xde709f2102306220921060314715629080e2fb77",
			name:    "All lowercase ETH address",
			token:   "ETH",
			isValid: true,
		},
		{
			address: "0XDE709F2102306220921060314715629080E2FB77",
			name:    "All uppercase ETH address",
			token:   "ETH",
			isValid: true,
		},
		{
			address: "0x52908400098527886e0F7030069857D2E4169EE7",
			name:    "Invalid checksum",
			token:   "ETH",
			isValid: false,
		},
		{
			address: "0x27b1FDB04752BBC536007A920D24ACB045561c26",
			name:    "Valid checksum",
			token:   "ETH",
			isValid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bool := ValidateAddress(tc.address, tc.token)
			assertCorrectValue(t, bool, tc.isValid)
		})
	}

}

func assertCorrectValue[T any](t testing.TB, got, want T) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
