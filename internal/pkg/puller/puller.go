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

	if fetchErr != nil && !errors.Is(fetchErr, git.NoErrAlreadyUpToDate) {
		err = fmt.Errorf("failed to fetch repository: %w", fetchErr)
		return
	}

	changes = true

	return
}
