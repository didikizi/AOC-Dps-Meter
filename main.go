package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"aocdpsmetr/internal/app"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// Create an instance of the app structure
	aocApp := app.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "AOC DPS Meter",
		Width:  1400,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        aocApp.Startup,
		OnDomReady:       aocApp.DomReady,
		OnBeforeClose:    aocApp.BeforeClose,
		OnShutdown:       aocApp.Shutdown,
		Debug: options.Debug{
			OpenInspectorOnStartup: true,
		},
		Bind: []interface{}{
			aocApp,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
