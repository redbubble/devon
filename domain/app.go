package domain

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var defaultSourceCodeBase = os.Getenv("HOME") + "/src"

const configFileName = "devon.conf.yaml"

type App struct {
	Name      string
	SourceDir string
}

type Config struct {
	Modes map[string]Mode
}

type Mode struct {
	Name         string
	Command      []string
	Dependencies map[string]string
}

func (a *App) Config() (Config, error) {
	// TODO: Memoize?
	var err error
	var bytes []byte
	var config Config

	configPath, err := a.configPath()

	if viper.IsSet("verbose") {
		fmt.Printf("Reading %s application metadata from: %s\n", a.Name, configPath)
	}

	if err != nil {
		return Config{}, err
	}

	bytes, err = ioutil.ReadFile(configPath)

	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func (a *App) sourceDir() (string, error) {
	if a.SourceDir != "" {
		return a.SourceDir, nil
	}

	path := filepath.Join(defaultSourceCodeBase, a.Name)

	if isDirectory(path) {
		return path, nil
	}

	return "", fmt.Errorf("Couldn't find a source code directory for '%s'.")
}

func (a *App) configPath() (string, error) {
	sourceDir, err := a.sourceDir()

	if err != nil {
		return "", err
	}

	return filepath.Join(sourceDir, configFileName), nil
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
