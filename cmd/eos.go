package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
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
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)

	rootCmd.AddCommand(newLoginCmd())
	rootCmd.AddCommand(newServerCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
