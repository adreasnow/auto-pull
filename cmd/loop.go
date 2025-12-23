package main

import (
	"context"
	"time"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
)

func loop(app *menuet.Application, ctx context.Context) {
	err := config.LoadConfig(ctx, app)
	if err != nil {
		app.Alert(menuet.Alert{
			MessageText:     "Failed to load config",
			InformativeText: err.Error(),
		})
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("failed to load config")
	}

	tickTime := time.Second * time.Duration(config.Config.RefreshSeconds)
	ticker := time.NewTicker(tickTime)

	// start a check immediately
	// must be in goroutine to be non-blocking
	go func() {
		tickNow <- struct{}{}
	}()

	for {
		select {
		case <-ticker.C:
			tick(ctx, app)
			ticker.Reset(tickTime)

		case <-tickNow:
			tick(ctx, app)
			ticker.Reset(tickTime)
		}
	}
}
