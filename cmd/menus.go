package main

import (
	"context"
	"os/exec"
	"strings"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
)

func menus(ctx context.Context, app *menuet.Application) (items []menuet.MenuItem) {
	return []menuet.MenuItem{
		{
			Type:    menuet.Regular,
			Text:    "Check Now",
			Clicked: checkNow,
		},
		{
			Type: menuet.Separator,
		},
		{
			Type:    menuet.Regular,
			Text:    "Show Directories",
			Clicked: func() { showDirectories(ctx, app) },
		},
		{
			Type:    menuet.Regular,
			Text:    "Show Logs",
			Clicked: func() { showLogs(ctx, app) },
		},
	}
}

func showDirectories(ctx context.Context, app *menuet.Application) {
	err := config.LoadConfig(ctx, app)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to load config")
		app.SetMenuState(&menuet.MenuState{Image: warningIcon})
		app.Alert(menuet.Alert{
			MessageText:     "Failed to load config",
			InformativeText: err.Error(),
		})
		return
	}

	builder := strings.Builder{}
	for _, dir := range config.Config.Directories {
		builder.WriteString("• " + dir + "\n") //nolint:errcheck
	}

	app.Alert(menuet.Alert{
		MessageText:     "Registered Directories",
		InformativeText: builder.String(),
	})
}

func showLogs(ctx context.Context, app *menuet.Application) {
	cmd := exec.Command("open", logFile)
	if err := cmd.Run(); err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to open logs")
		app.SetMenuState(&menuet.MenuState{Image: warningIcon})
		app.Alert(menuet.Alert{
			MessageText:     "Failed to open logs",
			InformativeText: err.Error(),
		})
		return
	}
}

func checkNow() {
	tickNow <- struct{}{}
}
