package cmd

import (
	// "fmt"

	"github.com/spf13/cobra"

	"github.com/redbubble/devon/domain"
)

var stopCmd = &cobra.Command{
	Use:   "stop [application] [flags]",
	Short: "Stop your chosen application",
	Long: `Stop your chosen application

* If <application> is set, devon will start it and its dependencies.
* If <application> is unset, devon will attempt to figure out which application
  it should start, based on the current working directory.
`,
	Run: func(cmd *cobra.Command, args []string) {
		var appName string
		var err error

		appName, err = getAppName(args)
		bail(err)

		app, err := domain.NewApp(appName, mode)
		bail(err)

		err = app.Stop()
		bail(err)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

	stopCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "development", "The mode to run in, e.g. 'development' or 'dependency'. Default: development")
}
