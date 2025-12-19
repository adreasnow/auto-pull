package main

import (
	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/adreasnow/auto-pull/internal/pkg/puller"
	"github.com/caseymrm/menuet"
)

func tick() {
	menuet.App().SetMenuState(&menuet.MenuState{Title: "🔍"})

	cfg, err := config.LoadConfig()
	if err != nil {
		menuet.App().SetMenuState(&menuet.MenuState{Title: "❌"})
		menuet.App().Notification(menuet.Notification{
			Title:    "Failed to load config",
			Subtitle: "",
			Message:  err.Error(),
		})
		return
	}

	success := true
	for _, dir := range cfg.Directories {
		changes, err := puller.Pull(dir)
		if err != nil {
			menuet.App().Notification(menuet.Notification{
				Title:    "Failed to fetch",
				Subtitle: dir,
				Message:  err.Error(),
			})
			success = false
			menuet.App().SetMenuState(&menuet.MenuState{Title: "❌"})
			continue
		}
		if changes {
			menuet.App().Notification(menuet.Notification{
				Title:                        "Changes detected",
				Subtitle:                     dir,
				RemoveFromNotificationCenter: true,
			})
		}
	}
	if success {
		menuet.App().SetMenuState(&menuet.MenuState{Title: "✅"})
	}
}
