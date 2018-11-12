package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "eos",
	Short: "Unified CLI to operate EOS clusters",
}

func Execute() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	viper.SetConfigFile(home + "/.eos/config.yml")
	viper.ReadInConfig()

	rootCmd.AddCommand(
		newLoginCmd(),
		newServerCmd(),
		newClusterCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
