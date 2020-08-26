package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var printVersion bool
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "devon start [application] [flags]",
	Short: "A tool to help you dev on your stuff",
	Long: `A tool to help you dev on your stuff

Devon starts applications for development, along with any other applications
they depend on.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if printVersion {
			versionCmd()
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(devonVersion string) {
	version = devonVersion

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devon.yaml)")

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Print all the informations!")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.Flags().BoolVar(&printVersion, "version", false, "Print the current version and exit")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".devon" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".devon")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func versionCmd() {
	fmt.Printf("devon v%s\n", version)
}
