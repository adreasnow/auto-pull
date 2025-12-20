package main

import (
	"context"
	"time"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
)

func loop(app *menuet.Application, ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(config.Config.RefreshSeconds))
	for {
		select {
		case <-ticker.C:
			tick(ctx, app)

		case <-tickNow:
			tick(ctx, app)
		}
	}
}
