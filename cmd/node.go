package cmd

import (
	"fmt"
	"os"

	"github.com/miguelaco/eos/common"
	"github.com/miguelaco/eos/config"
	"github.com/miguelaco/eos/mesos"
	"github.com/spf13/cobra"
)

func newNodeCmd() (cc *cobra.Command) {
	cc = &cobra.Command{
		Use:   "node",
		Short: "Manage EOS cluster nodes.",
	}

	cc.AddCommand(
		newNodeListCmd(),
	)

	return
}

func newNodeListCmd() (cac *cobra.Command) {
	cac = &cobra.Command{
		Use:   "list",
		Short: "List cluster nodes.",
		Run: func(cmd *cobra.Command, args []string) {
			cluster := config.GetAttachedCluster()
			if !cluster.Active {
				fmt.Println("No attached cluster")
				os.Exit(2)
			}

			fmt.Println("Cluster", cluster.Addr, "node list")

			nodes, err := mesos.Nodes(cluster)
			if err != nil {
				fmt.Println("Node list error:", err)
				os.Exit(3)
			}

			table := common.NewTable(os.Stdout, []string{"HOSTNAME", "ID", "TYPE"})
			for _, node := range nodes {
				table.Append([]string{node.Hostname, node.Id, node.Type})
			}

			table.Render()
		},
	}

	return
}
