package config

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/adreasnow/auto-pull/internal/github"
	"github.com/caseymrm/menuet"
	"github.com/rs/zerolog"
	"github.com/zalando/go-keyring"
)

var (
	service = "auto-pull-github-token"
	user    = os.Getenv("USER")

	ErrCancelledFlow = errors.New("user cancelled token entry")
)

func TokenFlow(ctx context.Context, app *menuet.Application) (token string, err error) {
	token, err = loadToken(ctx)
	if token != "" && err == nil {
		return
	}

	token, err = promptForToken(ctx, app)
	if err != nil {
		return // nolint:errcheck
	}

	err = saveToken(ctx, token)
	if err != nil {
		return // nolint:errcheck
	}

	err = github.CheckScopes(ctx, token)
	if err != nil {
		return // nolint:errcheck
	}

	return
}

func promptForToken(ctx context.Context, app *menuet.Application) (token string, err error) {
	zerolog.Ctx(ctx).Info().Msg("prompting user for token")
	response := app.Alert(menuet.Alert{
		MessageText:     "Token not found in keychain, please enter a GitHub token",
		InformativeText: "The GitHub token must have contents: read permissions",

		Inputs:  []string{"Token"},
		Buttons: []string{"Save Token", "Cancel"},
	})

	switch response.Button {
	case 0:
		zerolog.Ctx(ctx).Info().Msg("user provided token")
		token = response.Inputs[0]

	case 1:
		zerolog.Ctx(ctx).Info().Msg("user cancelled auth flow")
		err = ErrCancelledFlow
		return
	}
	return
}

func saveToken(ctx context.Context, token string) (err error) {
	if token == "" {
		return errors.New("token cannot be empty")
	}

	err = keyring.Set(service, user, token)
	if err != nil {
		err = fmt.Errorf("failed to save token to keychain: %w", err)
	}

	zerolog.Ctx(ctx).Info().Msg("saved token to keychain")

	return
}

func loadToken(ctx context.Context) (token string, err error) {
	token, err = keyring.Get(service, user)
	if err != nil {
		err = fmt.Errorf("failed to load token from keychain: %w", err)
		return
	}

	zerolog.Ctx(ctx).Info().Msg("loaded token from keychain")

	err = github.CheckScopes(ctx, token)
	if err != nil {
		err = fmt.Errorf("failed to validate presence of correct scopes: %w", err)
		return
	}

	return
}
