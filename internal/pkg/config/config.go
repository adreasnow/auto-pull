package config

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Directories    []string `yaml:"directories"`
	RefreshSeconds int      `yaml:"refreshSeconds"`
}

func LoadConfig() (cfg *Config, err error) {
	dirPath := path.Join(os.Getenv("HOME"), ".config/auto-pull")
	fileSystem := os.DirFS(dirPath)

	body, err := loadConfigfile(fileSystem)
	if err != nil {
		return
	}

	cfg, err = parseConfig(body)
	if err != nil {
		return
	}

	return
}

func loadConfigfile(fileSystem fs.FS) (body []byte, err error) {
	var file fs.File
	file, err = fileSystem.Open("config.yaml")
	if err != nil {
		err = fmt.Errorf("failed to read file from fs: %w", err)
		return
	}

	defer file.Close()

	body, err = io.ReadAll(file)
	if err != nil {
		err = fmt.Errorf("failed to read file content: %w", err)
		return
	}

	return
}

func parseConfig(data []byte) (cfg *Config, err error) {
	cfg = &Config{}
	if err = yaml.Unmarshal(data, cfg); err != nil {
		yaml.FormatError(err, true, true)
		err = fmt.Errorf("failed to unmarshal config file: %w", err)
		return
	}

	cfg.cleanTildeDirs()

	return
}

func (c *Config) cleanTildeDirs() {
	for i, dir := range c.Directories {
		if dir[0] == '~' {
			c.Directories[i] = strings.Replace(dir, "~", os.Getenv("HOME"), 1)
		}
	}
}
