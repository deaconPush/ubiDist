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
}

func NewETHAccount(masterKey *bip32.Key, tokenName string) (*ETHAccount, error) {
	privateKey, err := crypto.ToECDSA(masterKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to convert master key to ECDSA: %w", err)
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	return &ETHAccount{
		tokenName: tokenName,
		publicKey: publicKey,
	}, nil

}

func (account *ETHAccount) GetAddress() string {
	return crypto.PubkeyToAddress(*account.publicKey).Hex()
}

func (account *ETHAccount) RetrieveBalance(network string) (string, error) {
	if network == "" {
		network = defaultNetwork
	}
	provider := providers[network]

	balance, err := GetBalance(provider, account.GetAddress())
	if err != nil {
		return "", fmt.Errorf("error retrieving balance: %v", err)
	}

	return balance, nil
}

func (account *ETHAccount) GetTokenName() string {
	return account.tokenName
}
