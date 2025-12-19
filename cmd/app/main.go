package main

import (
	"os"
	"time"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
)

func loop(cfg *config.Config) {
	ticker := time.NewTicker(time.Second * time.Duration(cfg.RefreshSeconds))
	for range ticker.C {
		tick()
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		menuet.App().Alert(menuet.Alert{
			MessageText:     "Failed to load config",
			InformativeText: err.Error(),
		})
		os.Exit(1)
	}

	go loop(cfg)

	menuet.App().Label = "com.github.adreasnow.auto-pull"

	menuet.App().Children = menus

	menuet.App().SetMenuState(&menuet.MenuState{Title: "✅"})
	menuet.App().RunApplication()
}
