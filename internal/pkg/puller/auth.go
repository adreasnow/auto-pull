package puller

import (
	"errors"
	"fmt"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

var (
	ErrNoGithubToken = errors.New("githubToken config or GITHUB_TOKEN environment variable not set")
)

func (d *directory) setupAuth() (err error) {
	d.remote, err = d.repo.Remote("origin")
	if err != nil {
		err = fmt.Errorf("failed to get remote: %w", err)
		return
	}

	d.upstream = d.remote.Config().Name

	if len(d.remote.Config().URLs) == 0 {
		err = errors.New("no remote URLs found")
		return
	}

	tp, err := transport.NewEndpoint(d.remote.Config().URLs[0])
	if err != nil {
		err = fmt.Errorf("failed to create transport endpoint: %w", err)
		return
	}

	switch tp.Scheme {
	case "https":
		d.auth = &http.BasicAuth{
			Username: "x-access-token",
			Password: config.Config.GitHubToken,
		}

	default:
		err = fmt.Errorf("unsupported scheme: %s", tp.Scheme)
		return
	}

	return
}
