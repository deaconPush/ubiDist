# Pear Wallet

This project is a personally developed wallet application and is not intended for commercial use. Its primary purpose is to consolidate knowledge and learn about blockchain technology and cryptocurrencies.

## **Getting Started**

To start building or running the application, you need to have **Wails** and **Go** installed on your system. Follow the [official Wails documentation](https://wails.io/docs/gettingstarted/installation) for installation instructions.

Once Wails is installed, run the following command in your terminal to check for any missing dependencies:  
```bash
wails doctor
```
## Build and Run

To build the application you can run the following command at the wallet folder level: 

```bash
wails build -clean && ./build/bin/wallet
```

To work with the Hardhat network, you need to run the Hardhat node in a separate terminal. Instructions for this process can be found in the README inside the Hardhat folder.


## Testing hardhat RPC endpoints
You can test hardhat RPC endpoints by replacing main.go content and doing a manual build for the application:

```go
package main

import (
	"fmt"
	"wallet/internal/currencies"
)

func main() {
    // GETTING ACCOUNT ETH BALANCE
	balance, err := currencies.GetBalance("http://localhost:4545", "")
	if err != nil {
		println("Error:", err.Error())
	}
	etherBalance, err := currencies.HexToEther(balance)
	if err != nil {
		println("Error:", err.Error())
	}
	fmt.Println("balance is: ", etherBalance)
}
```

```go
package main

import (
	"log"
	"math/big"
	"wallet/internal/currencies"

	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
    // ETH TRANSACTION
    // You can retrieve the private keys from hardhat console output after starting the node
	privateKey, _ := crypto.HexToECDSA("<hex-private-key>")
	from := ""
	to := ""
	value := new(big.Int)
	ethToWeiMultiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
	value.Mul(big.NewInt(1000), ethToWeiMultiplier)
	signedTx, err := currencies.ProcessTransaction("http://localhost:8545", from, to, value, privateKey)
	if err != nil {
		log.Fatal("Error processing transaction:", err)
	}
	log.Println("Transaction sent:", signedTx)
}
```


You can build and run the application manually with the following command:
```bash
go build -o wallet && ./wallet
```


The operation will be registered by hardhat's node console:

![alt text](hardhat-output.png)