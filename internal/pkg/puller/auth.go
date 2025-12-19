package puller

import (
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/transport"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

func setupAuth(repo *git.Repository) (auth transport.AuthMethod, remote *git.Remote, err error) {
	remote, err = repo.Remote("origin")
	if err != nil {
		err = fmt.Errorf("failed to get remote: %w", err)
		return
	}

	cfg := remote.Config()

	tp, err := transport.NewEndpoint(cfg.URLs[0])
	if err != nil {
		err = fmt.Errorf("failed to create transport endpoint: %w", err)
		return
	}

	switch tp.URL.Scheme {
	case "ssh":
		auth, err = ssh.DefaultAuthBuilder("git")
		if err != nil {
			err = fmt.Errorf("failed to create SSH auth method: %w", err)
			return
		}

	case "https", "http":
		token, found := os.LookupEnv("GITHUB_TOKEN")
		if !found {
			err = fmt.Errorf("GITHUB_TOKEN environment variable not set")
			return
		}

		auth = &http.BasicAuth{
			Username: "x-access-token",
			Password: token,
		}

	default:
		err = fmt.Errorf("unsupported scheme: %s", tp.URL.Scheme)
		return
	}

	return
}
