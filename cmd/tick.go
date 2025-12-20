package main

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/adreasnow/auto-pull/internal/pkg/puller"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
)

func tick(ctx context.Context) {
	menuet.App().SetMenuState(&menuet.MenuState{Title: "💭"})

	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to load config")
		menuet.App().SetMenuState(&menuet.MenuState{Title: "⚠️"})
		menuet.App().Notification(menuet.Notification{
			Title:    "Failed to load config",
			Subtitle: "",
			Message:  err.Error(),
		})
		return
	}

	wg := sync.WaitGroup{}

	var success atomic.Bool
	success.Store(true)
	for _, dir := range cfg.Directories {
		wg.Go(func() {
			status := checkDir(ctx, cfg, dir)
			success.Store(status && success.Load())
		})
	}

	wg.Wait()

	if success.Load() {
		zerolog.Ctx(ctx).Info().Msg("successfully checked all directories")
		menuet.App().SetMenuState(&menuet.MenuState{Title: "👍"})
	}
}

func checkDir(ctx context.Context, cfg *config.Config, dir string) (success bool) {
	zerolog.Ctx(ctx).Info().Str("dir", dir).Msg("checking directory")
	changes, err := puller.Pull(cfg, dir)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("dir", dir).Msg("failed to fetch")
		menuet.App().SetMenuState(&menuet.MenuState{Title: "⚠️"})
		menuet.App().Notification(menuet.Notification{
			Title:    "Failed to fetch",
			Subtitle: dir,
			Message:  err.Error(),
		})
		return
	}
	if changes {
		zerolog.Ctx(ctx).Info().Str("dir", dir).Msg("changes detected")
		menuet.App().Notification(menuet.Notification{
			Title:    "Changes detected and pulled",
			Subtitle: dir,
		})
		success = true
		return
	}
	success = true
	zerolog.Ctx(ctx).Info().Str("dir", dir).Msg("no changes detected")
	return
}
