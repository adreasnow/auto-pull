package puller

import (
	"errors"
	"fmt"
	"os"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
)

func setupAuth(repo *git.Repository) (auth transport.AuthMethod, remote *git.Remote, err error) {
	remote, err = repo.Remote("origin")
	if err != nil {
		err = fmt.Errorf("failed to get remote: %w", err)
		return
	}

	remoteCfg := remote.Config()

	tp, err := transport.NewEndpoint(remoteCfg.URLs[0])
	if err != nil {
		err = fmt.Errorf("failed to create transport endpoint: %w", err)
		return
	}

	switch tp.Scheme {
	case "ssh":
		err = errors.New("ssh protocol not supported")
		return

	case "http":
		err = errors.New("http protocol not supported")
		return

	case "https":
		token, found := os.LookupEnv("GITHUB_TOKEN")
		if !found {
			token = config.Config.GithubToken
			if token == "" {
				err = fmt.Errorf("githubToken config or GITHUB_TOKEN environment variable not set")
				return
			}
		}

		auth = &http.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}

	default:
		err = fmt.Errorf("unsupported scheme: %s", tp.Scheme)
		return
	}

	return
}
