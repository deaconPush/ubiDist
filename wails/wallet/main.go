package main

import (
	"embed"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	//app := NewApp()
	getBalance("http://127.0.0.1:8545/", "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

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
