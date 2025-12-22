package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog"
)

func CheckScopes(ctx context.Context, token string) (err error) {
	if strings.HasPrefix(token, "github_pat_") {
		zerolog.Ctx(ctx).Info().
			Msg("cannot validate scopes for fine-grained pat")
		return
	}

	scopes, err := queryScopes(token)
	if err != nil {
		err = fmt.Errorf("failed to get scopes: %w", err)
		return
	}

	if scopes == "" {
		err = errors.New("no scopes found")
		return
	}

	fmt.Println(scopes)

	zerolog.Ctx(ctx).Info().
		Str("scopes", scopes).
		Msg("validated token scopes")

	return
}

func queryScopes(token string) (scopes string, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com", nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return
	}

	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to send request: %w", err)
		return
	}

	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return
	}

	scopes = resp.Header.Get("x-oauth-scopes")

	return
}
