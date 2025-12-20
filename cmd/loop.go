package main

import (
	"context"
	"time"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
)

func loop(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(config.Config.RefreshSeconds))
	for {
		select {
		case <-ticker.C:
			tick(ctx)

		case <-tickNow:
			tick(ctx)
		}
	}
}
