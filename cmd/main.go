package main

import (
	"context"
	"io"
	"os"
	"path"

	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	warningIcon = "warning.png"
	successIcon = "sun.png"
	pullingIcon = "cloud.png"
	bundleName  = "com.github.adreasnow.auto-pull"
	logFile     = path.Join(os.Getenv("HOME"), "Library", "Logs", bundleName, "logs.log")

	tickNow = make(chan struct{})
)

func main() {
	ctx := startLogger()

	app := initApp(ctx)

	go loop(app, ctx)

	app.RunApplication()
}

func startLogger() context.Context {
	logFile := lumberjack.Logger{
		Filename:   logFile,
		MaxBackups: 5,
		MaxSize:    2,
	}

	return zerolog.New(
		io.MultiWriter(
			zerolog.ConsoleWriter{Out: os.Stderr},
			zerolog.ConsoleWriter{Out: &logFile, NoColor: true}),
	).With().Timestamp().Logger().WithContext(context.Background())
}

func initApp(ctx context.Context) *menuet.Application {
	app := menuet.App()
	app.Label = bundleName
	app.Children = func() []menuet.MenuItem { return menus(ctx, app) }
	app.SetMenuState(&menuet.MenuState{Image: successIcon})
	app.NotificationResponder = func(id, response string) {}

	return app
}
