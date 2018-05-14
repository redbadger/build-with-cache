package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/redbadger/build-with-cache/constants"
	"github.com/redbadger/build-with-cache/root"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	file    string
	tag     string
)

var rootCmd = &cobra.Command{
	Use:   "build-with-cache PATH|URL|-",
	Short: "Docker build with layer caching via registry",
	Long: `
A cli command written in Go that uses a Docker registry to store layer caches in order to speed up build times. Useful in CI pipelines.

The tool parses the Dockerfile for the stage targets and attempts to pull respective images from the specified registry. Any images it finds are used as layer caches for the docker build. Updated images for each stage back are pushed back to the registry ready for the next build.
`,
	Version: constants.Version,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = root.Run(args[0], file, tag)
		return
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.build-with-cache.yaml)")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "Dockerfile", "Name of the Dockerfile (Default is 'PATH/Dockerfile')")
	rootCmd.PersistentFlags().StringVarP(&tag, "tag", "t", "", "Name and optionally a tag in the ‘registry/name:tag’ format")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".build-with-cache")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
