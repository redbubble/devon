package domain

import (
	"fmt"
	"os"
	"path/filepath"
)

var defaultSourceCodeBase = os.Getenv("HOME") + "/src"
const devonConfigFile = "devon.conf.yaml"

type App struct {
	Name string
	SourceDir string
}

func (a *App) DevonConfigPath() (string, error) {
	sourceDir, err := a.sourceDir()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", sourceDir, devonConfigFile), nil
}

func (a *App) sourceDir() (string, error) {
	if a.SourceDir != "" {
		return a.SourceDir, nil
	}

	path := filepath.Join(defaultSourceCodeBase, a.Name)

	if isDirectory(path) {
		return path, nil
	}

	// Golang is opinionated about where source code should live.
	gopath, gopathSet := os.LookupEnv("GOPATH")

	if gopathSet {
		path = filepath.Join(gopath, "src", "github.com", "redbubble", a.Name)

		if isDirectory(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("Couldn't find a source code directory for '%s'.")
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
