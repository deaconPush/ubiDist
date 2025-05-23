package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
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

func (a *ETHAccount) EstimateGas(to, value string) (string, error) {
	cliCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	valueWei, err := EtherToWei(value)
	if err != nil {
		return "", fmt.Errorf("error parsing ether transaction value: %v", err)
	}

	gasEstimate, err := a.client.EstimateGas(cliCtx, a.GetAddress(), to, valueWei)
	if err != nil {
		return "", fmt.Errorf("error estimating gas: %v", err)
	}

	gasPrice, err := a.client.GetGasPrice(cliCtx)
	if err != nil {
		return "", fmt.Errorf("error retrieving gas price %v", gasPrice)
	}

	gasEstimateBig := new(big.Int).SetUint64(gasEstimate)
	totalGasCostWei := new(big.Int).Mul(gasEstimateBig, gasPrice)
	weiFloat := new(big.Float).SetInt(totalGasCostWei)
	etherValue := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
	return etherValue.Text('f', 18), nil
}

func (a *ETHAccount) SendTransaction(to, value string, privateKey *ecdsa.PrivateKey) (string, error) {
	cliCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	weiValue, err := EtherToWei(value)
	if err != nil {
		return "", fmt.Errorf("error parsing ether value into wei: %v", err)
	}

	transactionHash, err := a.client.ProcessTransaction(cliCtx, a.GetAddress(), to, weiValue, privateKey)
	if err != nil {
		return "", fmt.Errorf("error procesing %s transaction %v", a.tokenName, err)
	}

	return transactionHash, nil
}

func (a *ETHAccount) GetTokenName() string {
	return a.tokenName
}

func (a *ETHAccount) ChangeProvider(provider string) {
	a.client.SetProvider(provider)
}
