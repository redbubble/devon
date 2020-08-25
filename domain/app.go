package domain

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

var DefaultSourceCodeBase = os.Getenv("HOME") + "/src"

const configFileName = "devon.conf.yaml"

type App struct {
	Name      string
	SourceDir string
	Config    Config
	Mode      Mode
}

type Config struct {
	Modes map[string]Mode
}

type Mode struct {
	Name         string
	Command      []string
	Dependencies map[string]string
}

func NewApp(name string, modeName string) (App, error) {
	var err error

	sourceDir := defaultSourceDir(name)

	if !isDirectory(sourceDir) {
		return App{}, fmt.Errorf("Couldn't find a source code directory for '%s'.", name)
	}

	if err != nil {
		return App{}, err
	}

	config, err := readConfig(name, sourceDir)

	if err != nil {
		return App{}, err
	}

	mode, ok := config.Modes[modeName]

	if !ok {
		return App{},
			fmt.Errorf("Couldn't find a mode called '%s' in the config for %s\n", modeName, name)
	}

	return App{
		Name:      name,
		SourceDir: sourceDir,
		Config:    config,
		Mode:      mode,
	}, nil
}

func readConfig(appName string, sourceDir string) (Config, error) {
	var err error
	var bytes []byte
	var config Config

	path := filepath.Join(sourceDir, configFileName)

	if viper.IsSet("verbose") {
		fmt.Printf("Reading %s application metadata from: %s\n", appName, path)
	}

	if err != nil {
		return Config{}, err
	}

	bytes, err = ioutil.ReadFile(path)

	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(bytes, &config)

	if err != nil {
		return Config{}, err
	}

	// Modifying properties of an object in a map is not allowed, so we have
	// to create a whole new object and assign that to the map.
	for name, mode := range config.Modes {
		newMode := Mode{
			Name:         name,
			Command:      mode.Command,
			Dependencies: mode.Dependencies,
		}

		config.Modes[name] = newMode
	}

	return config, nil
}

func defaultSourceDir(appName string) string {
	return filepath.Join(DefaultSourceCodeBase, appName)
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
