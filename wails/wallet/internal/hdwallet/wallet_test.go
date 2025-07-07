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

func TestWithHardhat(t *testing.T) {
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
	t.Run("Test wallet creation and operations", func(t *testing.T) {
		mnemonic := "test test test test test test test test test test test junk"
		password := "testpassword123#"
		to := "0x976EA74026E726554dB657fA54763abd0C3a0aa9"
		amountToSend := 530.0
		// Initialize in-memory wallet storage
		ws, err := hdwallet.NewWalletStorage(ctx, ":memory:")
		require.NoError(t, err)
		// Restore wallet using mnemonic and password
		wallet, err := hdwallet.RestoreWallet(ctx, password, mnemonic, ws)
		require.NoError(t, err)
		// Initialize ETH account
		err = wallet.Initialize([]string{"ETH"}, password, endpoint)
		require.NoError(t, err)
		// Check initial balance
		initialBalance, err := wallet.GetBalance("ETH", 0)
		require.NoError(t, err)
		require.Equal(t, float64(10000), initialBalance)
		// Estimate gas cost for the transaction
		gasPriceStr, err := wallet.EstimateGas("ETH", to, "530", 0)
		require.NoError(t, err)
		gasPrice, err := strconv.ParseFloat(gasPriceStr, 64)
		require.NoError(t, err)
		// Send the transaction
		_, err = wallet.SendTransaction("ETH", password, to, "530", 0)
		require.NoError(t, err)
		// Check final balance
		finalBalance, err := wallet.GetBalance("ETH", 0)
		require.NoError(t, err)
		// Validate balance with tolerance for float error
		expectedBalance := initialBalance - amountToSend - gasPrice
		const epsilon = 1e-6
		if math.Abs(finalBalance-expectedBalance) > epsilon {
			t.Errorf("expected balance after transfer error margin is too big")
		}
		// validate last transaction
		transactions, err := wallet.GetTransactions()
		latestTransaction := transactions[0]
		require.Equal(t, to, latestTransaction.Recipient)
		require.Equal(t, "530", latestTransaction.Value)
	})

}
