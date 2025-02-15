package eth

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
)

type ETHAccount struct {
	Token     string
	publicKey *ecdsa.PublicKey
}

func NewETHAccount(masterKey *bip32.Key, token string) (*ETHAccount, error) {
	privateKey, err := crypto.ToECDSA(masterKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert master key to ECDSA: %w", err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	return &ETHAccount{
		Token:     token,
		publicKey: publicKey,
	}, nil

}

func (account *ETHAccount) GetAddress() string {
	return crypto.PubkeyToAddress(*account.publicKey).Hex()
}

func (account *ETHAccount) RetrieveBalance(network string) (string, error) {
	balance, err := GetBalance(network, account.GetAddress())
	if err != nil {
		return "", fmt.Errorf("error retrieving balance: %v", err)
	}

	return balance, nil
}
