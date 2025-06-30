package eth

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"wallet/internal/utils"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
)

var providers = map[string]string{
	"hardhat": "http://localhost:8545",
}

const defaultNetwork = "hardhat"

type ETHMasterAccount struct {
	tokenName string
	client    *ethClient
	ctx       context.Context
	accountDB *AccountStorage
}

func NewETHAccount(ctx context.Context, masterKey *bip32.Key, tokenName string, db *sql.DB) (*ETHMasterAccount, error) {
	dbCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	accountDB, err := NewAccountStorage(db, dbCtx)
	defer cancel()

	if err != nil {
		return nil, fmt.Errorf("error initializing  %s account DB: %v", tokenName, err)
	}

	accountsExist, err := accountDB.AccountsExist(dbCtx)
	if err != nil {
		return nil, fmt.Errorf("error retrieving accounts from db: %v", err)
	}

	if !accountsExist {
		var ethAccounts []string

		for i := 0; i < 21; i++ {
			ethKey, err := utils.DeriveKeyForAccount(masterKey, "ETH", i)
			if err != nil {
				return nil, fmt.Errorf("error deriving %s account %d: %v", tokenName, i, err)
			}

			privateKey, err := crypto.ToECDSA(ethKey.Key)
			if err != nil {
				return nil, fmt.Errorf("failed to convert master key to ECDSA: %v", err)
			}

			ethAccounts = append(ethAccounts, crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		}
		err = accountDB.SaveAccounts(dbCtx, ethAccounts)
		if err != nil {
			return nil, fmt.Errorf("error saving %s accounts into the DB: %v", tokenName, err)
		}
	}

	client := NewEthClient(providers[defaultNetwork])
	return &ETHMasterAccount{
		tokenName: tokenName,
		client:    client,
		ctx:       ctx,
		accountDB: accountDB,
	}, nil

}

func (a *ETHMasterAccount) GetAddress(accountIndex int) (string, error) {
	dbCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	address, err := a.accountDB.GetAccountAddress(dbCtx, accountIndex)
	if err != nil {
		return "", fmt.Errorf("error retrieving account address from DB: %v", err)
	}

	return address, nil
}

func (a *ETHMasterAccount) GetAllAccounts() (map[int]string, error) {
	dbCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	return a.accountDB.GetAllAccounts(dbCtx)
}

func (a *ETHMasterAccount) RetrieveBalance(accountIndex int) (string, error) {
	cliCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()

	address, err := a.accountDB.GetAccountAddress(cliCtx, accountIndex)
	if err != nil {
		return "", fmt.Errorf("error retrieving account address from DB: %v", err)
	}

	balance, err := a.client.GetBalance(cliCtx, address)
	if err != nil {
		return "", fmt.Errorf("error retrieving balance: %v", err)
	}

	return balance, nil
}

func (a *ETHMasterAccount) EstimateGas(to, value string, accountIndex int) (string, error) {
	cliCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	valueWei, err := EtherToWei(value)
	if err != nil {
		return "", fmt.Errorf("error parsing ether transaction value: %v", err)
	}

	from, err := a.accountDB.GetAccountAddress(cliCtx, accountIndex)
	if err != nil {
		return "", fmt.Errorf("error retrieving account address from DB: %v", err)
	}

	gasEstimate, err := a.client.EstimateGas(cliCtx, from, to, valueWei)
	if err != nil {
		return "", fmt.Errorf("error estimating gas: %v", err)
	}

	gasPrice, err := a.client.GetGasPrice(cliCtx)
	if err != nil {
		return "", fmt.Errorf("error retrieving gas price %v", gasPrice)
	}

	totalCostEther := CalculateTotalGasCostInEther(gasEstimate, gasPrice)
	return totalCostEther, nil
}

func (a *ETHMasterAccount) SendTransaction(to, value string, masterKey *bip32.Key, accountIndex int) (string, error) {
	cliCtx, cancel := context.WithTimeout(a.ctx, 5*time.Second)
	defer cancel()
	weiValue, err := EtherToWei(value)
	if err != nil {
		return "", fmt.Errorf("error parsing ether value into wei: %v", err)
	}

	from, err := a.accountDB.GetAccountAddress(cliCtx, accountIndex)
	if err != nil {
		return "", fmt.Errorf("error retrieving account address from DB: %v", err)
	}

	ethKey, err := utils.DeriveKeyForAccount(masterKey, "ETH", accountIndex)
	if err != nil {
		return "", fmt.Errorf("error deriving eth key for index %d : %v", accountIndex, err)
	}

	privateKey, err := crypto.ToECDSA(ethKey.Key)
	if err != nil {
		return "", fmt.Errorf("failed to convert eth master key to ECDSA: %w", err)
	}

	transactionHash, err := a.client.ProcessTransaction(cliCtx, from, to, weiValue, privateKey)
	if err != nil {
		return "", fmt.Errorf("error procesing %s transaction %v", a.tokenName, err)
	}

	return transactionHash, nil
}

func (a *ETHMasterAccount) ChangeProvider(provider string) {
	a.client.SetProvider(provider)
}
