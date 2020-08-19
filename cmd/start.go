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

	"github.com/redbubble/devon/domain"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start your chosen application, along with its dependencies",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start called")

		// Figure out which app to start
		var appName string
		// var err error

		if len(args) > 0 {
			appName = args[0]
		} else {
			gitPath, err := currentGitRepo()
			appName = filepath.Base(gitPath)

			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		}

		// Read the devon config for that app
		app := domain.App{
			Name: appName,
		}

		fmt.Println(app.DevonConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")
	rootCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "development", "The mode to run in, e.g. 'development' or 'dependency'. Default: development")

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
