package utils

// Encryption implementation details from https://bruinsslot.jp/post/golang-crypto/
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"

	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/scrypt"
)

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("error generating entropy: %v", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("error generating mnemonic: %v", err)
	}
	return mnemonic, nil
}

func Encrypt(key, data []byte) ([]byte, error) {
	key, salt, err := deriveKey(key, nil)
	if err != nil {
		return nil, fmt.Errorf("error deriving key: %v", err)
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %v", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("error creating nonce: %v", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)

	return ciphertext, nil
}

func Decrypt(key, data []byte) ([]byte, error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := deriveKey(key, salt)
	if err != nil {
		return nil, fmt.Errorf("error deriving key: %v", err)
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("error creating cipher: %v", err)
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, fmt.Errorf("error creating GCM: %v", err)
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting data: %v", err)
	}

	return plainText, nil
}

func deriveKey(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, fmt.Errorf("error creating salt: %v", err)
		}
	}

	key, err := scrypt.Key(password, salt, 16384, 8, 1, 32)
	if err != nil {
		return nil, nil, fmt.Errorf("error deriving key: %v", err)
	}

	return key, salt, nil
}
