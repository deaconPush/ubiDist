package main

import (
	"log"
	"math/big"
	"wallet/internal/currencies"

	"github.com/ethereum/go-ethereum/crypto"
)

// import (
// 	"embed"
// )
// //go:embed all:frontend/dist
// //var assets embed.FS

func main() {
	// Create an instance of the app structure
	//app := NewApp()
	privateKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		log.Fatal("Error converting private key to ECDSA:", err)
	}
	from := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
	to := "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
	value := new(big.Int)
	ethToWeiMultiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil) // 10^18
	value.Mul(big.NewInt(1000), ethToWeiMultiplier)
	// modify the code to receive the value directly in eth, and then transform it to wei for the transaction
	signedTx, err := currencies.ProcessTransaction("http://localhost:8545", from, to, value, privateKey)
	if err != nil {
		log.Fatal("Error processing transaction:", err)
	}
	log.Println("Signed transaction:", signedTx)
	// Create application with options
	// err := wails.Run(&options.App{
	// 	Title:  "wallet",
	// 	Width:  1024,
	// 	Height: 768,
	// 	AssetServer: &assetserver.Options{
	// 		Assets: assets,
	// 	},
	// 	BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
	// 	OnStartup:        app.startup,
	// 	Bind: []interface{}{
	// 		app,
	// 	},
	// })

	// if err != nil {
	// 	println("Error:", err.Error())
	// }
}
