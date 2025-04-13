package eth

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
)

var providers = map[string]string{
	"hardhat": "http://localhost:8545",
}

const defaultNetwork = "hardhat"

type ETHAccount struct {
	tokenName string
	publicKey *ecdsa.PublicKey
	client    *ethClient
}

func NewETHAccount(masterKey *bip32.Key, tokenName string) (*ETHAccount, error) {
	privateKey, err := crypto.ToECDSA(masterKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert master key to ECDSA: %w", err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	client := NewEthClient(providers[defaultNetwork])
	return &ETHAccount{
		tokenName: tokenName,
		publicKey: publicKey,
		client:    client,
	}, nil

}

func (a *ETHAccount) GetAddress() string {
	return crypto.PubkeyToAddress(*a.publicKey).Hex()
}

func (a *ETHAccount) RetrieveBalance() (string, error) {
	balance, err := a.client.GetBalance(a.GetAddress())
	if err != nil {
		return "", fmt.Errorf("error retrieving balance: %v", err)
	}

	return balance, nil
}

func (a *ETHAccount) GetTokenName() string {
	return a.tokenName
}

func (a *ETHAccount) ChangeProvider(provider string) {
	a.client.SetProvider(provider)
}
