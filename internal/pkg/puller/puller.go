package puller

import (
	"errors"
	"fmt"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/go-git/go-git/v6"
)

func Pull(cfg *config.Config, path string) (changes bool, err error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		err = fmt.Errorf("failed to open repository for %s: %w", path, err)
		return
	}

	auth, remote, err := setupAuth(cfg, repo)
	if err != nil {
		err = fmt.Errorf("failed to setup auth for %s: %w", path, err)
		return
	}

	fetchErr := remote.Fetch(&git.FetchOptions{
		Auth:       auth,
		RemoteName: "origin",
	})

	if fetchErr != nil {
		if errors.Is(fetchErr, git.NoErrAlreadyUpToDate) {
			return
		}

		err = fmt.Errorf("failed to fetch repository for %s: %w", path, fetchErr)
		return
	}

	worktree, err := repo.Worktree()
	if err != nil {
		err = fmt.Errorf("failed to get worktree for %s: %w", path, err)
		return
	}

	pullErr := worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if pullErr != nil {
		if errors.Is(pullErr, git.NoErrAlreadyUpToDate) {
			err = fmt.Errorf("changes were detected but repo was already up to date on pull for %s: %w", path, pullErr)
			return
		}

		err = fmt.Errorf("failed to pull repository for %s: %w", path, pullErr)
		return
	}

	changes = true
	return
}
