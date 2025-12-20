package config

import (
	"embed"
	"fmt"
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

	cfg, err := parseConfig(data)
	require.NoError(t, err)

	assert.Len(t, cfg.Directories, 1)
	assert.Equal(t, cfg.Directories[0], "/test/path")
	assert.Equal(t, cfg.RefreshSeconds, 60)
}

func TestCleanTildeDirs(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		Directories: []string{"~/test/path"},
	}

	cfg.cleanTildeDirs()

	assert.Contains(t, cfg.Directories[0], "/Users")

	fmt.Println(cfg.Directories)
}
