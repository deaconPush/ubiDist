package main

import (
	"fmt"
	"log"
	"math/big"

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
	privateKey, err := crypto.HexToECDSA("")
	if err != nil {
		log.Fatal("Error converting private key to ECDSA:", err)
	}
	chainID := 1337
	from := ""
	to := ""
	signedTx := processTransaction("http://localhost:8545", from, to, big.NewInt(1000000000000000000), privateKey, int64(chainID))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("signedTx:", signedTx)
	// Get balance from receiver
	balance := getBalance("http://localhost:8545", to)
	fmt.Println("Balance after transaction:", balance)

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
