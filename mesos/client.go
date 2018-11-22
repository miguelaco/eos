package mesos

import (
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

const (
	basePath = "/mesos"
)

type Node struct {
	Id       string `json:id`
	Hostname string `json:hostname`
	Type     string
}

type Client struct {
	httpclient *common.HttpClient
	baseURL    string
}

func NewClient(cluster *config.Cluster) *Client {
	httpclient := common.NewHttpClient()
	httpclient.Token = cluster.Token

	baseURL := cluster.Addr + basePath
	return &Client{httpclient: httpclient, baseURL: baseURL}
}

func (c *Client) Nodes() (result map[string]Node, err error) {
	state, err := c.getMasterState()
	if err != nil {
		return
	}

	result = map[string]Node{}
	result[state.Hostname] = Node{Id: state.Id, Hostname: state.Hostname, Type: NodeTypeLeader}
	for _, slave := range state.Slaves {
		result[slave.Hostname] = Node{Id: slave.Id, Hostname: slave.Hostname, Type: NodeTypeSlave}
	}

	return
}

func (c *Client) Leader() (result Node, err error) {
	state, err := c.getStateSummary()
	if err != nil {
		return
	}

	result = Node{Id: state.Id, Hostname: state.Hostname, Type: NodeTypeLeader}
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

type stateSummary struct {
	Hostname string `json:hostname`
	Id       string `json:id`
	Slaves   []struct {
		Id       string `json:id`
		Hostname string `json:hostname`
	} `json:slaves`
}

func (c *Client) getStateSummary() (result stateSummary, err error) {
	res, err := c.get("/master/state-summary")
	if err != nil {
		return
	}

	defer res.Body.Close()

	result = stateSummary{}
	json.NewDecoder(res.Body).Decode(&result)
	return
}

func (c *Client) Verbose(verbose bool) {
	c.httpclient.Verbose(verbose)
}

func (c *Client) get(path string) (res *http.Response, err error) {
	url := c.baseURL + path
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
