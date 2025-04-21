package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

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
	ctx       context.Context
}

func NewETHAccount(ctx context.Context, masterKey *bip32.Key, tokenName string) (*ETHAccount, error) {
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
		ctx:       ctx,
	}, nil

}

func (a *ETHAccount) GetAddress() string {
	return crypto.PubkeyToAddress(*a.publicKey).Hex()
}

func (a *ETHAccount) RetrieveBalance() (string, error) {
	cliCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	balance, err := a.client.GetBalance(cliCtx, a.GetAddress())
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
