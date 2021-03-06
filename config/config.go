package config

import (
	"fmt"
	"io/ioutil"

	homedir "github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"
)

type Cluster struct {
	Addr   string `yaml:"addr,omitempty"`
	User   string `yaml:"user,omitempty"`
	Token  string `yaml:"token,omitempty"`
	Active bool   `yaml:active`
}

type Config struct {
	path     string
	clusters map[string]*Cluster
}

var c *Config

func init() {
	home, _ := homedir.Dir()
	c = &Config{}
	c.path = home + "/.eos"
	c.clusters = map[string]*Cluster{}
	c.Load()
}

func (cfg *Config) Load() error {
	data, err := ioutil.ReadFile(cfg.path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, &cfg.clusters); err != nil {
		return err
	}

	return nil
}

func Save() error {
	return c.Save()
}

func (cfg *Config) Save() error {
	y, err := yaml.Marshal(cfg.clusters)
	if err != nil {
		return err
	}

	if err = ioutil.WriteFile(cfg.path, y, 0644); err != nil {
		return err
	}

	return nil
}

func AddCluster(name string, cluster *Cluster) {
	c.AddCluster(name, cluster)
}

func (cfg *Config) AddCluster(name string, cluster *Cluster) {
	cfg.clusters[name] = cluster
}

func RemoveCluster(name string) error {
	return c.RemoveCluster(name)
}

func (cfg *Config) RemoveCluster(name string) error {
	if _, ok := cfg.clusters[name]; !ok {
		return fmt.Errorf("Cluster %s not found", name)
	}

	delete(cfg.clusters, name)

	return nil
}

func GetCluster(name string) *Cluster {
	return c.GetCluster(name)
}

func (cfg *Config) GetCluster(name string) *Cluster {
	return cfg.clusters[name]
}

func ListClusters() map[string]*Cluster {
	return c.ListClusters()
}

func (cfg *Config) ListClusters() map[string]*Cluster {
	return cfg.clusters
}

func AttachCluster(name string) error {
	return c.AttachCluster(name)
}

func (cfg *Config) AttachCluster(name string) error {
	if _, ok := cfg.clusters[name]; !ok {
		return fmt.Errorf("Cluster %s not found", name)
	}

	for n, cluster := range cfg.clusters {
		cluster.Active = false
		if n == name {
			cluster.Active = true
		}
	}

	return nil
}

func GetAttachedCluster() *Cluster {
	return c.GetAttachedCluster()
}

func (cfg *Config) GetAttachedCluster() *Cluster {
	for _, cluster := range cfg.clusters {
		if cluster.Active {
			return cluster
		}
	}

	return &Cluster{}
}
