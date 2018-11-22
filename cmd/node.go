package cmd

import (
	"fmt"
	"os"

	"github.com/miguelaco/eos/common"
	"github.com/miguelaco/eos/config"
	"github.com/miguelaco/eos/consul"
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

type NodeListCmd struct {
	*cobra.Command
	verbose bool
}

func newNodeListCmd() *cobra.Command {
	nlc := NodeListCmd{}

	nlc.Command = &cobra.Command{
		Use:   "list",
		Short: "List cluster nodes.",
		Run: func(cmd *cobra.Command, args []string) {
			cluster := config.GetAttachedCluster()
			if !cluster.Active {
				fmt.Println("No attached cluster")
				os.Exit(2)
			}

			fmt.Println("Cluster", cluster.Addr, "node list")

			consulClient := consul.NewClient(cluster)
			consulClient.Verbose(nlc.verbose)
			members, err := consulClient.Members()
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}

			mesosClient := mesos.NewClient(cluster)
			mesosClient.Verbose(nlc.verbose)
			leader, err := mesosClient.Leader()
			if err != nil {
				fmt.Println(err)
				os.Exit(3)
			}

			table := common.NewTable(os.Stdout, []string{"", "HOSTNAME", "IP", "STATUS", "TYPE"})
			for _, member := range members {
				if member.Addr == leader.Hostname {
					table.Append([]string{"*", member.Name, member.Addr, member.StatusText(), member.Type()})
				} else {
					table.Append([]string{"", member.Name, member.Addr, member.StatusText(), member.Type()})
				}
			}

			table.Render()
		},
	}

	nlc.Command.Flags().BoolVarP(&nlc.verbose, "verbose", "v", false, "Trace http requests")

	return nlc.Command
}
