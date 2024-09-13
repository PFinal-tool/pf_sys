package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "pf_sys",
		Width:  200,
		Height: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Menu:             nil,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.closeup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
	// systray.Run(onReady, nil)

}
