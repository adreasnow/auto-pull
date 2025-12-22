package puller

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v6"
	gitconfig "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullInvalidPath(t *testing.T) {
	t.Parallel()

	changes, pulled, commitMsg, repoName, err := Pull("/invalid/path/that/does/not/exist")

	require.ErrorContains(t, err, "failed to setup repository")
	assert.False(t, changes)
	assert.False(t, pulled)
	assert.Empty(t, commitMsg)
	assert.Empty(t, repoName)
}

func TestPullNoRemote(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	_, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	changes, pulled, commitMsg, repoName, err := Pull(tempDir)

	require.ErrorContains(t, err, "failed to setup repository")
	assert.False(t, changes)
	assert.False(t, pulled)
	assert.Empty(t, commitMsg)
	assert.Empty(t, repoName)
}

func TestPullNoAuth(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	require.NoError(t, err)

	t.Setenv("GITHUB_TOKEN", "")

	changes, pulled, commitMsg, repoName, err := Pull(tempDir)

	require.ErrorContains(t, err, "failed to setup repository")
	assert.False(t, changes)
	assert.False(t, pulled)
	assert.Empty(t, commitMsg)
	assert.Empty(t, repoName)
}

func TestFetchNoChanges(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	testFile := tempDir + "/test.txt"
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	remote, err := repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/nonexistent-repo.git"},
	})
	require.NoError(t, err)

	t.Setenv("GITHUB_TOKEN", "test-token")

	d := &directory{
		repoName: tempDir,
		repo:     repo,
		remote:   remote,
		upstream: "origin",
	}

	err = d.setupAuth()
	require.NoError(t, err)

	changes, err := d.fetch()
	require.Error(t, err)

	assert.False(t, changes, "no changes should be detected when fetch fails")
}

func TestPullDirtyWorkingTree(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	testFile := tempDir + "/test.txt"
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	require.NoError(t, err)

	err = os.WriteFile(testFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	t.Setenv("GITHUB_TOKEN", "test-token")

	changes, pulled, commitMsg, repoName, err := Pull(tempDir)

	require.ErrorContains(t, err, "failed to setup repository")
	assert.False(t, changes)
	assert.False(t, pulled)
	assert.Empty(t, commitMsg)

	// repoName is not set when setupRepo fails early
	assert.Empty(t, repoName)
}

func TestPull(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	testFile := tempDir + "/test.txt"
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	remote, err := repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	require.NoError(t, err)

	t.Setenv("GITHUB_TOKEN", "test-token")

	d := &directory{
		repoName: tempDir,
		repo:     repo,
		remote:   remote,
		upstream: "origin",
		worktree: worktree,
	}

	err = d.setupAuth()
	require.NoError(t, err)

	pulled, commitMsg, err := d.pull()

	// This will error due to network issues in unit tests
	// but we can verify the function structure

	assert.Error(t, err)

	assert.False(t, pulled, "should not pull when up to date")
	assert.Empty(t, commitMsg, "no commit message when not pulled")

}
