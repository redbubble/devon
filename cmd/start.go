/*
Copyright © 2020 Redbubble

*/
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

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start your chosen application, along with its dependencies",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Figure out which app to start
		var appName string
		var err error

		if len(args) > 0 {
			appName = args[0]
		} else {
			gitPath, err := currentGitRepo()
			bail(err)

			appName = filepath.Base(gitPath)
		}

		fmt.Printf("Starting %s in '%s' mode\n", appName, mode)

		app, err := domain.NewApp(appName, mode)
		bail(err)

		apps := make([]domain.App, 0, 1)
		apps, err = resolver.Add(apps, app)
		bail(err)

		if viper.IsSet("verbose") {
			fmt.Println("Devon will start these apps:")

			for i := 0; i < len(apps); i++ {
				fmt.Printf("* %s (in '%v' mode)\n", apps[i].Name, apps[i].Mode.Name)
			}
		}

	},
}

func bail(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")
	startCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "development", "The mode to run in, e.g. 'development' or 'dependency'. Default: development")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
