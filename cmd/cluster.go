package cmd

import (
	"fmt"
	"os"

	"github.com/miguelaco/eos/config"
	"github.com/spf13/cobra"
)

func newClusterCmd() (cc *cobra.Command) {
	cc = &cobra.Command{
		Use:   "cluster",
		Short: "Manage EOS clusters.",
	}

	cc.AddCommand(
		newClusterAddCmd(),
		newClusterListCmd(),
		newClusterAttachCmd(),
	)

	return
}

func newClusterAddCmd() (cac *cobra.Command) {
	cac = &cobra.Command{
		Use:   "add [name] [addr]",
		Short: "Add new cluster to config.",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			config.AddCluster(args[0], &config.Cluster{Addr: args[1]})
			config.AttachCluster(args[0])
			if err := config.Save(); err != nil {
				fmt.Println("Cannot save configuration:", err)
				os.Exit(3)
			}
		},
	}

	return
}

func newClusterListCmd() (cac *cobra.Command) {
	cac = &cobra.Command{
		Use:   "list",
		Short: "List configured clusters.",
		Run: func(cmd *cobra.Command, args []string) {
			list := config.ListClusters()
			fmt.Println(list)
		},
	}

	return
}

func newClusterAttachCmd() (cac *cobra.Command) {
	cac = &cobra.Command{
		Use:   "attach [name]",
		Short: "Attach to a given cluster.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config.AttachCluster(args[0])
			if err := config.Save(); err != nil {
				fmt.Println("Cannot save configuration:", err)
				os.Exit(3)
			}
		},
	}

	return
}
