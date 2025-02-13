package eth

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
)

type ETHAccount struct {
	token     string
	publicKey *ecdsa.PublicKey
}

func NewETHAccount(masterKey *bip32.Key, token string) (*ETHAccount, error) {
	privateKey, err := crypto.ToECDSA(masterKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert master key to ECDSA: %w", err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	return &ETHAccount{
		token:     token,
		publicKey: publicKey,
	}, nil

}

func (account *ETHAccount) GetAddress() string {
	return crypto.PubkeyToAddress(*account.publicKey).Hex()
}
