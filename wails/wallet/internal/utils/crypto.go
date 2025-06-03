package utils

// Encryption implementation details from https://bruinsslot.jp/post/golang-crypto/
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/sha3"
)

var TokenCoinTypes = map[string]int{
	"BTC": 0,  // Bitcoin
	"ETH": 60, // Ethereum
}

func GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(128)
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

func DeriveKeyForAccount(masterKey *bip32.Key, token string, accountIndex int) (*bip32.Key, error) {
	coinType, ok := TokenCoinTypes[token]
	if !ok {
		return nil, fmt.Errorf("token not found for account derivation")
	}

	derivationPath := fmt.Sprintf("m/44'/%d'/0'/0/%d", coinType, accountIndex)
	return deriveChildKey(masterKey, derivationPath)
}

func deriveChildKey(masterKey *bip32.Key, path string) (*bip32.Key, error) {
	indices, err := parseDerivationPath(path)
	if err != nil {
		return nil, fmt.Errorf("error parsing derivation path: %v", err)
	}
	key := masterKey
	for _, index := range indices {
		key, err = key.NewChildKey(index)
		if err != nil {
			return nil, fmt.Errorf("error deriving child key: %v", err)
		}
	}
	return key, nil
}

func parseDerivationPath(path string) ([]uint32, error) {
	var indices []uint32
	var hardenedOffset uint32 = 0x80000000
	for _, part := range strings.Split(path, "/") {
		if part == "m" {
			continue
		}

		hardened := strings.HasSuffix(part, "'")
		if hardened {
			part = part[:len(part)-1]
		}

		parsedIndex, err := strconv.ParseUint(part, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing derivation path: %v", err)
		}
		index := uint32(parsedIndex)

		if hardened {
			index += hardenedOffset
		}

		indices = append(indices, index)
	}
	return indices, nil
}

func ValidateAddress(address, token string) bool {
	if token == "ETH" {
		return ValidateETHAddress(address)

	}
	return false
}

func ValidateETHAddress(address string) bool {
	ok := ValidateETHAddressFormat(address)
	if !ok {
		return false
	}

	cleanAddress := address[2:]
	if cleanAddress == strings.ToLower(cleanAddress) || cleanAddress == strings.ToUpper(cleanAddress) {
		return true
	}
	ok = ValidateETHChecksum(cleanAddress)

	return ok
}

func ValidateETHAddressFormat(address string) bool {
	ethAddressRegex := regexp.MustCompile(`^0[xX][0-9a-fA-F]{40}$`)
	return ethAddressRegex.MatchString(address)

}

func ValidateETHChecksum(address string) bool {
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(strings.ToLower(address)))
	digest := hasher.Sum(nil)
	digestHex := hex.EncodeToString(digest)
	var checksumAddress strings.Builder

	for i, c := range address {
		if c >= '0' && c <= '9' {
			checksumAddress.WriteRune(c)
			continue
		}

		hashBit, err := strconv.ParseUint(string(digestHex[i]), 16, 4)

		if err != nil {
			return false
		}

		if hashBit > 7 {
			checksumAddress.WriteRune(unicode.ToUpper(c))
		} else {
			checksumAddress.WriteRune(c)
		}
	}
	return address == checksumAddress.String()
}
