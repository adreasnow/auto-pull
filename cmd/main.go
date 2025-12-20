package main

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	warningIcon = "warning.png"
	successIcon = "sun.png"
	pullingIcon = "cloud.png"
	bundleName  = "com.github.adreasnow.auto-pull"

	tickNow = make(chan struct{})
)

func main() {
	ctx := context.Background()

	logFile := lumberjack.Logger{
		Filename:   path.Join(os.Getenv("HOME"), "Library", "Logs", bundleName, "logs.log"),
		MaxBackups: 5,
		MaxSize:    2,
	}

	ctx = zerolog.New(
		io.MultiWriter(
			zerolog.ConsoleWriter{Out: os.Stderr},
			&logFile),
	).With().Timestamp().Logger().WithContext(ctx)

	err := config.LoadConfig(ctx)
	if err != nil {
		menuet.App().Alert(menuet.Alert{
			MessageText:     "Failed to load config",
			InformativeText: err.Error(),
		})
		zerolog.Ctx(ctx).Fatal().Err(err).Msg("failed to load config")
	}

	go loop(ctx)

	menuet.App().Label = bundleName

	menuet.App().Children = func() []menuet.MenuItem { return menus(ctx) }

	menuet.App().SetMenuState(&menuet.MenuState{Image: successIcon})
	menuet.App().RunApplication()
}
