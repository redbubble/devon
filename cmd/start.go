package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/redbubble/devon/domain"
	"github.com/redbubble/devon/resolver"
)

var mode string
var skip []string

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start [application] [flags]",
	Short: "Start your chosen application, along with its dependencies",
	Long: `Start your chosen application, along with its dependencies

* If <application> is set, devon will start it and its dependencies.
* If <application> is unset, devon will attempt to figure out which application
  it should start, based on the current working directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Figure out which app to start
		var appName string
		var err error

		appName, err = getAppName(args)
		bail(err)

		fmt.Printf("Starting %s in '%s' mode\n", appName, mode)

		app, err := domain.NewApp(appName, mode)
		bail(err)

		apps := make([]domain.App, 0, 1)
		apps, err = resolver.Add(apps, app, skip)
		bail(err)

		if len(apps) == 0 {
			bail(fmt.Errorf("There are no apps to start. Consider checking %s/devon.conf.yaml to figure out what's going on.", app.SourceDir))
		}

		if viper.IsSet("verbose") {
			fmt.Println("Devon will start these apps:")

			printAppList(apps)
		}

		// TODO: Make this display some help for tmux -- not everyone knows and loves it
		// TODO: Consider not hard-coding the tmux session name
		c := exec.Command("tmux", "new-session", "-d", "-s", "devon", "devon --version | less")
		_, err = c.Output()

		bail(err)

		// Iterate through the list backwards. This gives the best
		// chance (though it doesn't guarantee) that dependencies will
		// be started before their dependents.
		startedApps := make([]domain.App, 0, len(apps))
		failedApps := make([]domain.App, 0, len(apps))

		for i := len(apps) - 1; i >= 0; i-- {
			err = apps[i].Start()

			if err == nil {
				startedApps = append(startedApps, apps[i])
			} else {
				fmt.Printf("WARNING: Could not start %s: %v\n", apps[i].Name, err)
				failedApps = append(failedApps, apps[i])
			}
		}

		if viper.IsSet("verbose") {
			fmt.Printf("These apps started successfully:\n")
			printAppList(startedApps)
			fmt.Println()
		}

		if len(failedApps) > 0 {
			fmt.Printf("Some apps failed to start:\n")
			printAppList(failedApps)
			// os.Exit(1)
		}

		fmt.Printf("Attaching to tmux session...")

		executable, err := exec.LookPath("tmux")
		bail(err)

		tmuxCmd := exec.Cmd{
			Path: executable,
			Args: []string{"tmux", "attach", "-t", "devon"},
			// Dir:    workingDir,
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		}

		err = tmuxCmd.Run()
		bail(err)

		for _, app := range startedApps {
			err = app.Stop()

			if err != nil {
				fmt.Printf("WARNING: Could not stop %s: %v\n", app.Name, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "development", "The mode to run in, e.g. 'development' or 'dependency'. Default: development")
	startCmd.PersistentFlags().StringSliceVarP(&skip, "skip", "s", []string{}, "Names of any apps you don't want to start. Can be specified multiple times.")
}

func getAppName(args []string) (string, error) {
	if len(args) > 0 {
		return args[0], nil
	}

	gitPath, err := currentGitRepo()

	if err != nil {
		return "", err
	}

	return filepath.Base(gitPath), nil
}

func bail(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		os.Exit(1)
	}
}

func currentGitRepo() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()

	if err != nil {
		return "", fmt.Errorf("We don't seem to be in a Git repo. Please provide the name of an application to start.")
	}

	trimmedOutput := strings.TrimSuffix(string(output), "\n")

	return trimmedOutput, nil
}

func printAppList(apps []domain.App) {
	for _, a := range apps {
		fmt.Printf("* %s (in '%s' mode)\n", a.Name, a.Mode.Name)
	}
}
