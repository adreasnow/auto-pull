package puller

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v6"
)

func Pull(path string) (changes bool, err error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		err = fmt.Errorf("failed to open repository: %w", err)
		return
	}

	auth, remote, err := setupAuth(repo)
	if err != nil {
		err = fmt.Errorf("failed to setup auth: %w", err)
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

		err = fmt.Errorf("failed to fetch repository: %w", fetchErr)
		return
	}

	changes = true

	worktree, err := repo.Worktree()
	if err != nil {
		err = fmt.Errorf("failed to get worktree: %w", err)
		return
	}

	pullErr := worktree.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if pullErr != nil {
		if errors.Is(pullErr, git.NoErrAlreadyUpToDate) {
			err = fmt.Errorf("changes were detected but repo was already up to date on pull: %w", pullErr)
			return
		}

		err = fmt.Errorf("failed to pull repository: %w", pullErr)
		return
	}

	return
}
