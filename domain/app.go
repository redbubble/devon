package domain

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"

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
	StartCommand []string `yaml:"start-command"`
	StopCommand  []string `yaml:"stop-command"`
	WorkingDir   string   `yaml:"working-dir"`
	Dependencies []Dependency
	Env          EnvVars
}

type Dependency struct {
	AppName  string `yaml:"name"`
	ModeName string `yaml:"mode"`
}

type EnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type EnvVars []EnvVar

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

func (a *App) Start() error {
	fmt.Println()
	fmt.Printf("----- Starting %s -----\n", a.Name)

	if len(a.Mode.StartCommand) == 0 {
		return fmt.Errorf("Don't know how to start %s in '%s' mode, because start-command is unset.", a.Name, a.Mode.Name)
	}

	return runCommand(a.Mode.StartCommand, a.Mode.WorkingDir, a.Mode.Env)
}

func (a *App) Stop() error {
	fmt.Println()
	fmt.Printf("----- Stopping %s -----\n", a.Name)

	if len(a.Mode.StopCommand) == 0 {
		return fmt.Errorf("Don't know how to stop %s in '%s' mode, because stop-command is unset.", a.Name, a.Mode.Name)
	}

	return runCommand(a.Mode.StopCommand, a.Mode.WorkingDir, a.Mode.Env)
}

func (env EnvVars) ToStrings() []string {
	strings := []string{}

	for _, variable := range env {
		strings = append(strings, variable.Name+"="+variable.Value)
	}

	return strings
}

func runCommand(command []string, workingDir string, env EnvVars) error {

	// In the case of tools like `make`, the executable will be on the PATH
	// and LookPath will find it.
	executable, err := exec.LookPath(command[0])

	// If not, the executable may be a script within the repo. LookPath will
	// only find that if we give it a complete path, including the
	// application directory.
	if err != nil {
		executable, err = exec.LookPath(filepath.Join(workingDir, command[0]))
	}

	if err != nil {
		return err
	}

	cmd := exec.Cmd{
		Path:   executable,
		Args:   command,
		Dir:    workingDir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    env.ToStrings(),
	}

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

		workingDir := filepath.Join(sourceDir, mode.WorkingDir)

		newMode := Mode{
			Name:         name,
			StartCommand: mode.StartCommand,
			StopCommand:  mode.StopCommand,
			Dependencies: mode.Dependencies,
			WorkingDir:   workingDir,
			Env:          mode.Env,
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
