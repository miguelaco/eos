package mesos

import (
	//	"fmt"
	"encoding/json"
	"errors"
	"net/http"

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

type Client struct {
	httpclient *common.HttpClient
	cluster    *config.Cluster
}

func NewClient(cluster *config.Cluster) *Client {
	httpclient := common.NewHttpClient()
	httpclient.Token = cluster.Token
	return &Client{httpclient: httpclient, cluster: cluster}
}

func (c *Client) Nodes() (result []Node, err error) {
	state, err := c.getMasterState()
	if err != nil {
		return
	}

	result = append(result, Node{Id: state.Id, Hostname: state.Hostname, Type: NodeTypeLeader})
	for _, slave := range state.Slaves {
		result = append(result, Node{Id: slave.Id, Hostname: slave.Hostname, Type: NodeTypeSlave})
	}

	return
}

type masterState struct {
	Hostname string `json:hostname`
	Id       string `json:id`
	Slaves   []struct {
		Id       string `json:id`
		Hostname string `json:hostname`
	} `json:slaves`
}

func (c *Client) getMasterState() (result masterState, err error) {
	res, err := c.get("/master/state")
	if err != nil {
		return
	}

	defer res.Body.Close()

	result = masterState{}
	json.NewDecoder(res.Body).Decode(&result)
	return
}

func (c *Client) get(path string) (res *http.Response, err error) {
	url := c.cluster.Addr + "/mesos" + path

	res, err = c.httpclient.Get(url)
	if err != nil {
		return
	}

	if res.StatusCode != http.StatusOK {
		err = errors.New(res.Status)
		return
	}

	return
}
