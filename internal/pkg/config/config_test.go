package config

import (
	"embed"
	"fmt"
	"strings"
	"testing"

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
}

func TestParseConfig(t *testing.T) {
	t.Parallel()

	data := []byte("directories:\n  - /test/path\nrefreshSeconds: 60")

	c := &config{}

	err := c.parseConfig(data)
	require.NoError(t, err)

	assert.Len(t, c.Directories, 1)
	assert.Equal(t, c.Directories[0], "/test/path")
	assert.Equal(t, c.RefreshSeconds, 60)
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

	fmt.Println(c.Directories)
}

func TestExpandTrailingWildCardDirs(t *testing.T) {
	t.Parallel()

	c := &config{
		Directories: []string{
			"/usr/*",
			"/test/dir",
		},
	}

	err := c.expandTrailingWildCardDirs()
	require.NoError(t, err)

	assert.Greater(t, len(c.Directories), 2)
	assert.Contains(t, c.Directories, "/usr/lib")
}

func TestCheckForGit(t *testing.T) {
	t.Parallel()

	c := &config{
		Directories: []string{
			"/usr/lib",
			"/Users/adrea/Downloads/dotfiles",
		},
	}

	c.checkForGit()

	assert.Len(t, c.Directories, 1)
	assert.NotContains(t, c.Directories, "/usr/lib")
}
