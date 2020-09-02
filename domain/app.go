package domain

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

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

	repoPath, err := getRepoPath(name)

	if err != nil {
		return App{}, err
	}

	config, err := readConfig(name, repoPath)

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
		SourceDir: repoPath,
		Config:    config,
		Mode:      mode,
	}, nil
}

func (a *App) Start() error {
	executable, err := a.executable()

	cmd := exec.Cmd{
		Path:   executable,
		Args:   a.Mode.Command,
		Dir:    a.SourceDir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Println()
	fmt.Printf("----- Starting %s -----\n", a.Name)

	if viper.IsSet("verbose") {
		fmt.Printf("Working directory: %s\n", cmd.Dir)
		fmt.Printf("Command: %v\n", cmd.Args)
		fmt.Println()
	}

	err = cmd.Start()

	if err != nil {
		return err
	}

	return cmd.Wait()
}

func (a *App) executable() (string, error) {
	command := a.Mode.Command

	// In the case of tools like `make`, the executable will be on the PATH
	// and LookPath will find it.
	executable, err := exec.LookPath(command[0])

	// If not, the executable may be a script within the repo. LookPath will
	// only find that if we give it a complete path, including the
	// application directory.
	if err != nil {
		executable, err = exec.LookPath(filepath.Join(a.SourceDir, command[0]))
	}

	if err != nil {
		return "", err
	}

	return executable, nil
}

func getRepoPath(appName string) (string, error) {
	var err error
	sourceDir := defaultSourceDir(appName)

	if !isDirectory(sourceDir) {
		repoUrlFormat := viper.GetString("repo-url-format")
		if viper.GetBool("clone-missing-repos") {
			repoUrl := fmt.Sprintf(repoUrlFormat, appName)

			fmt.Printf("INFO: Cloning %s into %s\n", repoUrl, sourceDir)

			cmd := exec.Command("git", "clone", repoUrl, sourceDir)
			out, err := cmd.CombinedOutput()

			if viper.IsSet("verbose") {
				for _, line := range strings.Split(string(out), "\n") {
					fmt.Printf("> %s\n", line)
				}
			}

			if err != nil {
				return "", fmt.Errorf("Couldn't find a source code directory for '%s', and couldn't `git clone` it.", appName)
			}
		} else {
			return "", fmt.Errorf("Couldn't find a source code directory for '%s'. Consider enabling --clone-missing-repos.", appName)
		}
	}

	if err != nil {
		return "", err
	}

	return sourceDir, nil
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
	return filepath.Join(viper.GetString("source-code-base-dir"), appName)
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
