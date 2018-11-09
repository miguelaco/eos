package cmd

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    jww "github.com/spf13/jwalterweatherman"
    homedir "github.com/mitchellh/go-homedir"
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