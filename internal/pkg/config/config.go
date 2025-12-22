package config

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/caseymrm/menuet"
	"github.com/goccy/go-yaml"
	"github.com/rs/zerolog"
)

var (
	Config = &config{}
)

var (
	ErrNoDirectories    = fmt.Errorf("no directories specified in config")
	ErrNoRefreshSeconds = fmt.Errorf("no refreshSeconds specified in config")
)

type config struct {
	Directories    []string `yaml:"directories"`
	RefreshSeconds int      `yaml:"refreshSeconds"`
	GitHubToken    string   `yaml:"-"`
}

func LoadConfig(ctx context.Context, app *menuet.Application) (err error) {
	dirPath := path.Join(os.Getenv("HOME"), ".config/auto-pull")
	fileSystem := os.DirFS(dirPath)

	body, err := loadConfigfile(fileSystem)
	if err != nil {
		err = fmt.Errorf("failed to load config file: %v", err)
		return
	}

	err = Config.parseConfig(body)
	if err != nil {
		err = fmt.Errorf("failed to parse config file: %w", err)
		return
	}

	if len(Config.Directories) == 0 {
		err = ErrNoDirectories
		return
	}

	if Config.RefreshSeconds == 0 {
		err = ErrNoRefreshSeconds
		return
	}

	if Config.GitHubToken == "" {
		Config.GitHubToken, err = TokenFlow(ctx, app)
		if err != nil {
			err = fmt.Errorf("failed to get github token: %w", err)
			return
		}
	}

	zerolog.Ctx(ctx).Info().
		Strs("directories", Config.Directories).
		Int("refreshSeconds", Config.RefreshSeconds).
		Msg("config loaded")

	return
}

func loadConfigfile(fileSystem fs.FS) (body []byte, err error) {
	var file fs.File
	file, err = fileSystem.Open("config.yaml")
	if err != nil {
		err = fmt.Errorf("failed to read file from fs: %w", err)
		return
	}

	defer file.Close() //nolint:errcheck

	body, err = io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("failed to read file content: %w", err)
		return
	}

	return
}

func (c *config) parseConfig(data []byte) (err error) {
	if err = yaml.Unmarshal(data, c); err != nil {
		yaml.FormatError(err, true, true)
		err = fmt.Errorf("failed to unmarshal config file: %w", err)
		return
	}

	c.cleanTildeDirs()
	err = c.expandTrailingWildCardDirs()
	if err != nil {
		err = fmt.Errorf("failed to expand trailing wild card directories: %w", err)
		return
	}
	c.checkForGit()

	return
}

func (c *config) cleanTildeDirs() {
	for i, dir := range c.Directories {
		if strings.HasPrefix(dir, "~") {
			c.Directories[i] = strings.Replace(dir, "~", os.Getenv("HOME"), 1)
		}
	}
}

func (c *config) expandTrailingWildCardDirs() (err error) {
	newDirectoryList := []string{}

	for _, dir := range c.Directories {
		if strings.HasSuffix(dir, "/*") {
			parentDir := dir[0 : len(dir)-2]

			files, err := os.ReadDir(parentDir)
			if err != nil {
				err = fmt.Errorf("failed to read directory: %w", err)
				return err
			}

			for _, f := range files {
				if f.IsDir() {
					newDirectoryList = append(newDirectoryList, filepath.Join(parentDir, f.Name()))
				}
			}
		} else {
			newDirectoryList = append(newDirectoryList, dir)
		}
	}

	c.Directories = newDirectoryList

	return
}

func (c *config) checkForGit() {
	newDirectoryList := []string{}

	for _, dir := range c.Directories {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			newDirectoryList = append(newDirectoryList, dir)
		}
	}

	c.Directories = newDirectoryList
}
