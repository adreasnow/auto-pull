package puller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/transport"
)

type directory struct {
	repoName      string
	upstream      string
	defaultBranch string

	repo     *git.Repository
	auth     transport.AuthMethod
	remote   *git.Remote
	worktree *git.Worktree
	head     *plumbing.Reference
}

func (d *directory) setupRepo() (err error) {
	d.repo, err = git.PlainOpen(d.repoName)
	if err != nil {
		err = fmt.Errorf("failed to open repository: %w", err)
		return
	}

	err = d.setupAuth()
	if err != nil {
		err = fmt.Errorf("failed to setup auth: %w", err)
		return
	}

	d.getRepoName()
	err = d.getDefaultBranch() //nolint:errcheck

	return
}

func (d *directory) getDefaultBranch() (err error) {
	refs, err := d.remote.List(&git.ListOptions{
		Auth: d.auth,
	})
	if err != nil {
		err = fmt.Errorf("failed to list remote references: %w", err)
		return
	}

	for _, ref := range refs {
		if ref.Name() == plumbing.HEAD {
			d.defaultBranch = ref.Target().Short()
			return
		}
	}

	err = errors.New("failed to find default branch")
	return
}

func (d *directory) getCommitMessage() (msg string, err error) {
	d.head, err = d.repo.Head()
	if err != nil {
		err = fmt.Errorf("failed to get head: %w", err)
		return
	}

	commit, err := d.repo.CommitObject(d.head.Hash())
	if err != nil {
		err = fmt.Errorf("failed to get commit: %w", err)
		return
	}

	msg = commit.Message
	return
}

func (d *directory) getRepoName() {
	cfg := d.remote.Config()
	d.repoName = strings.Replace(cfg.URLs[0], "https://github.com/", "", 1)
	d.repoName = strings.Replace(d.repoName, ".git", "", 1)
}

func (d *directory) checkBranchStatus() (isClean bool, err error) {
	d.worktree, err = d.repo.Worktree()
	if err != nil {
		err = fmt.Errorf("failed to get worktree %w", err)
		return
	}

	status, err := d.worktree.Status()
	if err != nil {
		err = fmt.Errorf("failed to get status: %w", err)
		return
	}

	if !status.IsClean() {
		return false, nil
	}

	for _, fileStatus := range status {
		if fileStatus.Worktree == git.Untracked {
			return false, nil
		}
	}

	return true, nil
}
