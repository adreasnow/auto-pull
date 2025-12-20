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

func tick(ctx context.Context, app *menuet.Application) {
	app.SetMenuState(&menuet.MenuState{Image: pullingIcon})

	err := config.LoadConfig(ctx)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Msg("failed to load config")
		app.SetMenuState(&menuet.MenuState{Image: warningIcon})
		app.Notification(menuet.Notification{
			Title:    "Failed to load config",
			Subtitle: "",
			Message:  err.Error(),

			Identifier: "config",
		})
		return
	}

	wg := sync.WaitGroup{}

	var success atomic.Bool
	success.Store(true)
	for _, dir := range config.Config.Directories {
		wg.Go(func() {
			status := checkDir(ctx, app, dir)
			success.Store(status && success.Load())
		})
	}

	wg.Wait()

	if success.Load() {
		zerolog.Ctx(ctx).Info().Msg("successfully checked all directories")
		app.SetMenuState(&menuet.MenuState{Image: successIcon})
	}
}

func checkDir(ctx context.Context, app *menuet.Application, dir string) (success bool) {
	zerolog.Ctx(ctx).Info().Str("dir", dir).Msg("checking directory")
	changes, err := puller.Pull(dir)
	if err != nil {
		zerolog.Ctx(ctx).Error().Err(err).Str("dir", dir).Msg("failed to fetch")
		app.SetMenuState(&menuet.MenuState{Image: warningIcon})
		app.Notification(menuet.Notification{
			Title:    "Failed to fetch",
			Subtitle: dir,
			Message:  err.Error(),

			Identifier: dir,
		})
		return
	}
	if changes {
		zerolog.Ctx(ctx).Info().Str("dir", dir).Msg("changes detected")
		app.Notification(menuet.Notification{
			Title:    "Changes detected and pulled",
			Subtitle: dir,

			Identifier: dir,
		})
		success = true
		return
	}
	success = true
	zerolog.Ctx(ctx).Info().Str("dir", dir).Msg("no changes detected")
	return
}
