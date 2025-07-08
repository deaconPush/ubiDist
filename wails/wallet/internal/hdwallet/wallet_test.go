package hdwallet_test

import (
	"context"
	"math"
	"strconv"
	"testing"
	"time"
	"wallet/internal/hdwallet"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestHDWallet(t *testing.T) {
	defaultMnemonic := "test test test test test test test test test test test junk"
	defaultPassword := "testpassword123#"
	defaultWrongPwd := "testpassword125#"
	firstAccountIndex := 0

	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../../",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"8545/tcp"},
		WaitingFor:   wait.ForLog("Started HTTP and WebSocket JSON-RPC server").WithStartupTimeout(120 * time.Second),
	}

	hardhatC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	require.NoError(t, err)
	endpoint, err := hardhatC.Endpoint(ctx, "")
	require.NoError(t, err)
	endpoint = "http://" + endpoint
	// Initialize in-memory wallet storage
	ws, err := hdwallet.NewWalletStorage(ctx, ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() {
		ws.Close()
		err = hardhatC.Terminate(ctx)
		if err != nil {
			t.Fatalf("error terminating hardhat container")
		}
	})
	t.Run("Test wallet restore and operations", func(t *testing.T) {
		to := "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
		transferAmountFloat := 530.0
		transferAmount := "530"
		expectedInitialBalance := float64(10000)
		tokens := []string{"ETH"}
		operationsToken := "ETH"
		// Restore wallet using mnemonic and password
		wallet, err := hdwallet.RestoreWallet(ctx, defaultPassword, defaultMnemonic, ws)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize(tokens, defaultPassword, endpoint)
		require.NoError(t, err)
		// Check initial balance
		initialBalance, err := wallet.GetBalance(operationsToken, firstAccountIndex)
		require.NoError(t, err)
		require.Equal(t, expectedInitialBalance, initialBalance)
		// Estimate gas cost for the transaction
		gasPriceStr, err := wallet.EstimateGas(operationsToken, to, transferAmount, firstAccountIndex)
		require.NoError(t, err)
		gasPrice, err := strconv.ParseFloat(gasPriceStr, 64)
		require.NoError(t, err)
		// Send the transaction
		_, err = wallet.SendTransaction(operationsToken, defaultPassword, to, transferAmount, firstAccountIndex)
		require.NoError(t, err)
		// Check final balance
		finalBalance, err := wallet.GetBalance(operationsToken, firstAccountIndex)
		require.NoError(t, err)
		// Validate balance with tolerance for float error
		expectedBalance := initialBalance - transferAmountFloat - gasPrice
		const epsilon = 1e-6
		if math.Abs(finalBalance-expectedBalance) > epsilon {
			t.Errorf("expected balance after transfer error margin is too big")
		}
		// validate last transaction
		transactions, err := wallet.GetTransactions()
		require.NoError(t, err)
		latestTransaction := transactions[0]
		require.Equal(t, to, latestTransaction.Recipient)
		require.Equal(t, transferAmount, latestTransaction.Value)
	})

	t.Run("Test wallet recovery and operations", func(t *testing.T) {
		var firstAcctPreviousBalance = 9469.999961
		var sixthAcctPreviousBalance = 10530
		tokens := []string{"ETH"}
		operationsToken := "ETH"
		firstAccountIndex := 0
		sixthAccountIndex := 6
		wallet, err := hdwallet.RecoverWallet(ctx, defaultPassword, ws)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize(tokens, defaultPassword, endpoint)
		require.NoError(t, err)
		// check balance of 0th and 6th accountRecoverWallet
		firstAccBalance, err := wallet.GetBalance(operationsToken, firstAccountIndex)
		require.NoError(t, err)
		sixthAcctBalance, err := wallet.GetBalance(operationsToken, sixthAccountIndex)
		require.NoError(t, err)
		require.Equal(t, firstAcctPreviousBalance, firstAccBalance)
		require.Equal(t, sixthAcctPreviousBalance, sixthAcctBalance)
	})

	t.Run("Test wallet restore with invalid mnemonic", func(t *testing.T) {
		mnemonic := "apple banana giraffe moon tiger salad potato elephant monkey window ocean"
		newWs, err := hdwallet.NewWalletStorage(ctx, ":memory:")
		require.NoError(t, err)
		_, err = hdwallet.RestoreWallet(ctx, defaultPassword, mnemonic, newWs)
		require.Error(t, err)
		newWs.Close()
	})

	t.Run("Test wallet recovery with wrong password", func(t *testing.T) {
		_, err := hdwallet.RecoverWallet(ctx, defaultWrongPwd, ws)
		require.Error(t, err)
	})

	t.Run("Test wallet initialization with unsupported password and token", func(t *testing.T) {
		tokens := []string{"ETH"}
		unsupportedTokens := []string{"BCH"}
		newWs, err := hdwallet.NewWalletStorage(ctx, ":memory:")
		require.NoError(t, err)
		wallet, err := hdwallet.RestoreWallet(ctx, defaultPassword, defaultMnemonic, newWs)
		require.NoError(t, err)
		// Initialize wallet with unsupported token
		err = wallet.Initialize(unsupportedTokens, defaultPassword, endpoint)
		require.Error(t, err)
		// Initialize wallet with invalid password
		err = wallet.Initialize(tokens, defaultWrongPwd, endpoint)
		require.Error(t, err)
		newWs.Close()
	})

	t.Run("Test wallet operations with invalid password and token", func(t *testing.T) {
		unsupportedToken := "BCH"
		tokens := []string{"ETH"}
		operationsToken := "ETH"
		senderAccountIndex := 0
		to := "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
		transferAmount := "530"
		wallet, err := hdwallet.RecoverWallet(ctx, defaultPassword, ws)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize(tokens, defaultPassword, endpoint)
		require.NoError(t, err)
		// Send transaction with invalid pwd
		_, err = wallet.SendTransaction(operationsToken, defaultWrongPwd, to, transferAmount, senderAccountIndex)
		require.Error(t, err)
		// Send transaction with invalid token
		_, err = wallet.SendTransaction(unsupportedToken, defaultWrongPwd, to, transferAmount, senderAccountIndex)
		require.Error(t, err)
		// Test getBalance with invalid token
		_, err = wallet.GetBalance(unsupportedToken, senderAccountIndex)
		require.Error(t, err)
	})

	t.Run("Test sending more tokens than allowed", func(t *testing.T) {
		var to = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"
		tokens := []string{"ETH"}
		operationsToken := "ETH"
		transferAmount := "15000"
		senderAccountIndex := 1
		wallet, err := hdwallet.RecoverWallet(ctx, defaultPassword, ws)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize(tokens, defaultPassword, endpoint)
		require.NoError(t, err)
		// Attempt to send transaction
		_, err = wallet.SendTransaction(operationsToken, defaultPassword, to, transferAmount, senderAccountIndex)
		require.Error(t, err)
	})

	t.Run("Test get all accounts method", func(t *testing.T) {
		expectedAccounts := map[int]string{
			0:  "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
			1:  "0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
			2:  "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC",
			3:  "0x90F79bf6EB2c4f870365E785982E1f101E93b906",
			4:  "0x15d34AAf54267DB7D7c367839AAf71A00a2C6A65",
			5:  "0x9965507D1a55bcC2695C58ba16FB37d819B0A4dc",
			6:  "0x976EA74026E726554dB657fA54763abd0C3a0aa9",
			7:  "0x14dC79964da2C08b23698B3D3cc7Ca32193d9955",
			8:  "0x23618e81E3f5cdF7f54C3d65f7FBc0aBf5B21E8f",
			9:  "0xa0Ee7A142d267C1f36714E4a8F75612F20a79720",
			10: "0xBcd4042DE499D14e55001CcbB24a551F3b954096",
		}
		unsupportedToken := "BCH"
		tokens := []string{"ETH"}
		operationsToken := "ETH"
		wallet, err := hdwallet.RecoverWallet(ctx, defaultPassword, ws)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize(tokens, defaultPassword, endpoint)
		require.NoError(t, err)
		accounts, err := wallet.GetAllAccounts(operationsToken)
		require.NoError(t, err)
		require.Equal(t, expectedAccounts, accounts)
		_, err = wallet.GetAllAccounts(unsupportedToken)
		require.Error(t, err)
	})

	t.Run("Test wallet creation and operations", func(t *testing.T) {
		tokens := []string{"ETH"}
		operationsToken := "ETH"
		firstAccountIndex := 0
		expectedBalance := float64(0)
		newWs, err := hdwallet.NewWalletStorage(ctx, ":memory:")
		require.NoError(t, err)
		wallet, mnemonic, err := hdwallet.CreateWallet(ctx, defaultPassword, newWs)
		require.NotEmpty(t, mnemonic)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize(tokens, defaultPassword, endpoint)
		require.NoError(t, err)
		// Retrieve balance
		balance, err := wallet.GetBalance(operationsToken, firstAccountIndex)
		require.NoError(t, err)
		require.Equal(t, expectedBalance, balance)
		newWs.Close()
	})
}
