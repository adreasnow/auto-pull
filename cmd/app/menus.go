package main

import (
	"fmt"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
)

func menus() (items []menuet.MenuItem) {
	items = []menuet.MenuItem{}

	showRepos := menuet.MenuItem{
		Type:    menuet.Regular,
		Text:    "Show Directories",
		Clicked: showDirectories,
	}

	items = append(items,
		showRepos,
	)

	return
}

func showDirectories() {
	cfg, err := config.LoadConfig()
	if err != nil {
		menuet.App().SetMenuState(&menuet.MenuState{Title: "❌"})
		menuet.App().Alert(menuet.Alert{
			MessageText:     "Failed to load config",
			InformativeText: err.Error(),
		})
		return
	}

	dirString := ""
	for _, dir := range cfg.Directories {
		dirString += fmt.Sprintf("• %s\n", dir)
	}

	menuet.App().Alert(menuet.Alert{
		MessageText:     "Registered Directories",
		InformativeText: dirString,
	})

}
