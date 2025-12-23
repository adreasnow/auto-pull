package puller

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v6"
)

func Pull(path string) (changes bool, pulled bool, commitMsg string, repoName string, err error) {
	dir := &directory{
		repoName: path,
	}

	err = dir.setupRepo()
	if err != nil {
		if errors.Is(err, git.ErrRemoteNotFound) {
			err = nil
			repoName = "local-only"
		} else {
			err = fmt.Errorf("failed to setup repository %s: %w", path, err)
		}

		return
	}

	repoName = dir.repoName

	changes, err = dir.fetch()
	if err != nil {
		err = fmt.Errorf("failed to fetch repository for %s: %w", path, err)
		return
	}

	isClean, err := dir.checkBranchStatus()
	if err != nil {
		err = fmt.Errorf("failed to check branch status: %w", err)
		return
	}

	if isClean {
		pulled, commitMsg, err = dir.pull()
		if err != nil {
			err = fmt.Errorf("failed to pull repository for %s: %w", path, err)
			return
		}
	}

	return
}

func (d *directory) fetch() (changes bool, err error) {
	fetchErr := d.remote.Fetch(&git.FetchOptions{
		Auth:       d.auth,
		RemoteName: d.upstream,
	})

	if fetchErr != nil {
		if errors.Is(fetchErr, git.NoErrAlreadyUpToDate) {
			return
		}

		err = fmt.Errorf("failed to fetch repository: %w", fetchErr)
		return
	}

	changes = true

	return
}

func (d *directory) pull() (pulled bool, commitMsg string, err error) {
	pullErr := d.worktree.Pull(&git.PullOptions{
		RemoteName: d.upstream,
		Auth:       d.auth,
	})

	if pullErr != nil {
		if errors.Is(pullErr, git.NoErrAlreadyUpToDate) ||
			errors.Is(pullErr, git.ErrWorktreeNotClean) ||
			errors.Is(pullErr, git.ErrUnstagedChanges) {
			return
		}

		err = fmt.Errorf("failed to pull repository: %w", pullErr)
		return
	}

	commitMsg, err = d.getCommitMessage()
	if err != nil {
		err = fmt.Errorf("failed to get commit message: %w", err)
		return
	}

	pulled = true

	return
}
