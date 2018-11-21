package mesos

import (
	//	"fmt"
	"encoding/json"

	"github.com/miguelaco/eos/common"
	"github.com/miguelaco/eos/config"
)

const (
	NodeTypeSlave  = "slave"
	NodeTypeMaster = "master"
	NodeTypeLeader = "master (leader)"
)

type Node struct {
	Id       string `json:id`
	Hostname string `json:hostname`
	Type     string
}

type stateSummary struct {
	Hostname string `json:hostname`
	Id       string `json:id`
	Slaves   []struct {
		Id       string `json:id`
		Hostname string `json:hostname`
	} `json:slaves`
}

func Nodes(cluster *config.Cluster) (result []Node, err error) {
	state, err := getStateSummary(cluster)
	if err != nil {
		return
	}

	result = append(result, Node{Id: state.Id, Hostname: state.Hostname, Type: NodeTypeLeader})
	for _, slave := range state.Slaves {
		result = append(result, Node{Id: slave.Id, Hostname: slave.Hostname, Type: NodeTypeSlave})
	}

	return
}

func getStateSummary(cluster *config.Cluster) (result stateSummary, err error) {
	url := cluster.Addr + "/mesos/master/state"
	c := common.NewHttpClient()
	c.Token = cluster.Token

	res, err := c.Get(url)
	if err != nil {
		return
	}

	defer res.Body.Close()

	result = stateSummary{}
	json.NewDecoder(res.Body).Decode(&result)
	return
}
