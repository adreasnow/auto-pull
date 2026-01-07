package config

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-git/go-git/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed config.yaml
var testData embed.FS

func TestLoad(t *testing.T) {
	t.Parallel()

	body, err := loadConfigfile(testData)
	require.NoError(t, err)

	assert.Contains(t, string(body), "- /test/path")
	assert.Contains(t, string(body), "refreshSeconds: 60")
	assert.Contains(t, string(body), "notifications:")
	assert.Contains(t, string(body), "  failed: true")
	assert.Contains(t, string(body), "  fetchedNoPull: true")
	assert.Contains(t, string(body), "  success: true")
}

func TestParseConfig(t *testing.T) {
	t.Parallel()

	t.Run("directories", func(t *testing.T) {
		t.Parallel()

		// Create a temporary directory with a git repo
		tempDir := t.TempDir()
		_, err := git.PlainInit(tempDir, false)
		require.NoError(t, err)

		data := fmt.Appendf([]byte{}, "directories:\n  - %s", tempDir)

		c := &config{}

		err = c.parseConfig(data)
		require.NoError(t, err)

		assert.Len(t, c.Directories, 1)
		assert.Equal(t, c.Directories[0], tempDir)
	})

	t.Run("notifications", func(t *testing.T) {
		t.Parallel()

		body, err := loadConfigfile(testData)
		require.NoError(t, err)

		c := &config{}

		err = c.parseConfig(body)
		require.NoError(t, err)

		assert.Equal(t, c.RefreshSeconds, 60)
		assert.True(t, c.Notifications.Failed)
		assert.True(t, c.Notifications.FetchedNoPull)
		assert.True(t, c.Notifications.Pulled)
	})
}

func TestCleanTildeDirs(t *testing.T) {
	t.Parallel()

	c := &config{
		Directories: []string{"~/test/path"},
	}

	c.cleanTildeDirs()

	splitLine := strings.Split(c.Directories[0], "/")
	require.Greater(t, len(splitLine), 1)

	assert.Contains(t, []string{"Users", "home"}, splitLine[1])
	assert.NotContains(t, c.Directories[0], "~")
}

func TestExpandTrailingWildCardDirs(t *testing.T) {
	t.Parallel()

	tempParent := t.TempDir()
	subDir1 := tempParent + "/subdir1"
	subDir2 := tempParent + "/subdir2"
	err := os.Mkdir(subDir1, 0755)
	require.NoError(t, err)
	err = os.Mkdir(subDir2, 0755)
	require.NoError(t, err)

	c := &config{
		Directories: []string{
			tempParent + "/*",
			"/test/dir",
		},
	}

	err = c.expandTrailingWildCardDirs()
	require.NoError(t, err)

	assert.Greater(t, len(c.Directories), 2)
	assert.Contains(t, c.Directories, subDir1)
	assert.Contains(t, c.Directories, subDir2)
	assert.Contains(t, c.Directories, "/test/dir")
}

func TestCheckForGit(t *testing.T) {
	t.Parallel()

	tempDirWithGit := t.TempDir()
	tempDirWithoutGit := t.TempDir()

	_, err := git.PlainInit(tempDirWithGit, false)
	require.NoError(t, err)

	c := &config{
		Directories: []string{
			tempDirWithoutGit,
			tempDirWithGit,
		},
	}

	c.checkForGit()

	assert.Len(t, c.Directories, 1)
	assert.NotContains(t, c.Directories, tempDirWithoutGit)
	assert.Contains(t, c.Directories, tempDirWithGit)
}
