package main

import (
	"context"
	"strings"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
)

func menus(ctx context.Context, app *menuet.Application) (items []menuet.MenuItem) {
	items = []menuet.MenuItem{}

	showRepos := menuet.MenuItem{
		Type:    menuet.Regular,
		Text:    "Show Directories",
		Clicked: func() { showDirectories(ctx, app) },
	}

	checkNow := menuet.MenuItem{
		Type:    menuet.Regular,
		Text:    "Check Now",
		Clicked: checkNow,
	}

	items = append(items,
		showRepos,
		checkNow,
	)

	return
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

func checkNow() {
	tickNow <- struct{}{}
}
