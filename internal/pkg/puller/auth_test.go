package puller

import (
	"testing"

	"github.com/adreasnow/auto-pull/internal/pkg/config"
	"github.com/go-git/go-git/v6"
	gitconfig "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupAuthHTTPS(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/test/repo.git"},
	})
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	t.Run("success", func(t *testing.T) {
		testToken := "test-token-123"
		t.Setenv("GITHUB_TOKEN", testToken)

		err = d.setupAuth()
		require.NoError(t, err)

		assert.NotNil(t, d.remote)
		assert.Equal(t, "origin", d.upstream)
		assert.NotNil(t, d.auth)

		basicAuth, ok := d.auth.(*http.BasicAuth)
		require.True(t, ok, "auth should be BasicAuth")
		assert.Equal(t, "x-access-token", basicAuth.Username)
		assert.Equal(t, testToken, basicAuth.Password)
	})

	t.Run("success from config", func(t *testing.T) {
		t.Setenv("GITHUB_TOKEN", "")
		config.Config.GithubToken = "config-token-456"

		err = d.setupAuth()
		require.NoError(t, err)

		assert.NotNil(t, d.remote)
		assert.Equal(t, "origin", d.upstream)
		assert.NotNil(t, d.auth)

		basicAuth, ok := d.auth.(*http.BasicAuth)
		require.True(t, ok, "auth should be BasicAuth")
		assert.Equal(t, "x-access-token", basicAuth.Username)
		assert.Equal(t, "config-token-456", basicAuth.Password)
	})

	t.Run("failure no token", func(t *testing.T) {
		t.Setenv("GITHUB_TOKEN", "")
		config.Config.GithubToken = ""

		err = d.setupAuth()
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrNoGithubToken)
	})
}

func TestSetupAuthNoRemote(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	err = d.setupAuth()
	require.ErrorContains(t, err, "failed to get remote")
}

func TestSetupAuthUnsupportedScheme(t *testing.T) {
	tempDir := t.TempDir()

	repo, err := git.PlainInit(tempDir, false)
	require.NoError(t, err)

	_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
		Name: "origin",
		URLs: []string{"git@github.com:test/repo.git"},
	})
	require.NoError(t, err)

	d := &directory{
		repoName: tempDir,
		repo:     repo,
	}

	t.Setenv("GITHUB_TOKEN", "test-token-123")

	err = d.setupAuth()
	require.ErrorContains(t, err, "unsupported scheme")
}
