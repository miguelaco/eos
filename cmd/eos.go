package cmd

import (
    "fmt"
    "os"
    "log"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
    Use:   "eos",
    Short: "Unified CLI to operate EOS clusters",
}

func Execute() {
    rootCmd.AddCommand(newLoginCmd())

    viper.SetConfigFile("/home/majimenez/.eos/config.yml")

    if err := viper.ReadInConfig(); err != nil {
        fmt.Println("Can't read config:", err)
        os.Exit(1)
    }

    log.Println(viper.AllKeys())

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}