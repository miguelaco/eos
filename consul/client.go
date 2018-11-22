package consul

import (
	"encoding/json"
	"errors"
	"net/http"
	"sort"

	"github.com/miguelaco/eos/common"
	"github.com/miguelaco/eos/config"
)

const (
	basePath = "/consul"
)

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

func (c *Client) Members() (result MemberSlice, err error) {
	res, err := c.get("/v1/agent/members")
	if err != nil {
		return
	}

	defer res.Body.Close()

	result = MemberSlice{}
	json.NewDecoder(res.Body).Decode(&result)
	sort.Sort(result)

	return
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
