package main

import (
	"context"
	"strings"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
)

func menus(ctx context.Context) (items []menuet.MenuItem) {
	items = []menuet.MenuItem{}

	showRepos := menuet.MenuItem{
		Type:    menuet.Regular,
		Text:    "Show Directories",
		Clicked: func() { showDirectories(ctx) },
	}

	items = append(items,
		showRepos,
	)

	return
}

func showDirectories(ctx context.Context) {
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to load config")
		menuet.App().SetMenuState(&menuet.MenuState{Title: "❌"})
		menuet.App().Alert(menuet.Alert{
			MessageText:     "Failed to load config",
			InformativeText: err.Error(),
		})
		return
	}

	builder := strings.Builder{}
	for _, dir := range cfg.Directories {
		builder.WriteString("• " + dir + "\n") //nolint:errcheck
	}

	menuet.App().Alert(menuet.Alert{
		MessageText:     "Registered Directories",
		InformativeText: builder.String(),
	})
}
