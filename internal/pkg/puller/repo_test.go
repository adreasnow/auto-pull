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

func TestGetRepoName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "GitHub HTTPS URL with .git",
			url:      "https://github.com/adreasnow/auto-pull.git",
			expected: "adreasnow/auto-pull",
		},
		{
			name:     "GitHub HTTPS URL without .git",
			url:      "https://github.com/adreasnow/auto-pull",
			expected: "adreasnow/auto-pull",
		},
		{
			name:     "GitHub URL with organization",
			url:      "https://github.com/golang/go.git",
			expected: "golang/go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()

			repo, err := git.PlainInit(tempDir, false)
			require.NoError(t, err)

			remote, err := repo.CreateRemote(&gitconfig.RemoteConfig{
				Name: "origin",
				URLs: []string{tt.url},
			})
			require.NoError(t, err)

			d := &directory{
				repoName: tempDir,
				repo:     repo,
				remote:   remote,
			}

			d.getRepoName()

			assert.Equal(t, tt.expected, d.repoName)
		})
	}
}

func TestCheckBranchStatusClean(t *testing.T) {
	t.Parallel()

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

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	isClean, err := d.checkBranchStatus()
	require.NoError(t, err)

	assert.True(t, isClean, "repository should be clean")
	assert.NotNil(t, d.worktree)
}

func TestCheckBranchStatusWithModifications(t *testing.T) {
	t.Parallel()

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

	err = os.WriteFile(testFile, []byte("modified content"), 0644)
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	isClean, err := d.checkBranchStatus()
	require.NoError(t, err)

	assert.False(t, isClean, "repository should not be clean with modifications")
}

func TestCheckBranchStatusWithUntrackedFiles(t *testing.T) {
	t.Parallel()

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

	untrackedFile := tempDir + "/untracked.txt"
	err = os.WriteFile(untrackedFile, []byte("untracked content"), 0644)
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	isClean, err := d.checkBranchStatus()
	require.NoError(t, err)

	assert.False(t, isClean, "repository should not be clean with untracked files")
}

func TestGetCommitMessage(t *testing.T) {
	t.Parallel()

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

	expectedMessage := "Test commit message\n\nThis is a detailed commit message."
	_, err = worktree.Commit(expectedMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	msg, err := d.getCommitMessage()
	require.NoError(t, err)

	assert.Equal(t, expectedMessage, msg)
	assert.NotNil(t, d.head)
}

func TestGetCommitMessageNoCommits(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	_, err = d.getCommitMessage()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get head")
}

func TestSetupRepoInvalidPath(t *testing.T) {
	t.Parallel()

	d := &directory{
		repoName: "/invalid/path/that/does/not/exist",
	}

	err := d.setupRepo()
	require.ErrorContains(t, err, "failed to open repository")
}
