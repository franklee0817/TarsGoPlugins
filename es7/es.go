package es7

import (
	"errors"
	"sync"

	"github.com/franklee0817/TarsGoPlugins/manager"
	"github.com/olivere/elastic/v7"
	"gopkg.in/yaml.v3"
)

const (
	PluginType = "clients"
	PluginName = "es7"
)

type clients struct {
	mux     *sync.RWMutex
	Clients map[string]*elastic.Client
}

var Clients *clients

type config struct {
	URL         string `yaml:"url"`          // es连接地址
	User        string `yaml:"user"`         // es连接用户名
	Password    string `yaml:"password"`     // es连接密码
	Sniff       bool   `yaml:"sniff"`        // SetSniff的值
	HealthCheck bool   `yaml:"health_check"` // 是否开启健康检查,默认关闭
}

var cfgs map[string]config

func init() {
	once := new(sync.Once)
	once.Do(func() {
		cfgs = make(map[string]config)
		if Clients == nil {
			Clients = new(clients)
			Clients.mux = new(sync.RWMutex)
			Clients.Clients = make(map[string]*elastic.Client)
			manager.Register(Clients)
		}
	})
}

func (c *clients) Type() string {
	return PluginType
}

func (c *clients) Name() string {
	return PluginName
}

func (c *clients) Setup(yamlCfg *yaml.Node) error {
	err := yamlCfg.Decode(&cfgs)
	if err != nil {
		return err
	}
	c.mux.Lock()
	defer c.mux.Unlock()
	for instName, cfg := range cfgs {
		c.Clients[instName], err = elastic.NewClient(
			elastic.SetURL(cfg.URL),
			elastic.SetSniff(cfg.Sniff),
			elastic.SetHealthcheck(cfg.HealthCheck),
			elastic.SetBasicAuth(cfg.User, cfg.Password))
		if err != nil {
			return err
		}
	}

	return nil
}

// GetClient 获取es Client
func GetClient(instName string) (*elastic.Client, error) {
	Clients.mux.RLock()
	defer Clients.mux.RUnlock()
	if clt, ok := Clients.Clients[instName]; ok {
		return clt, nil
	}

	return nil, errors.New("no es client instance found for " + instName)
}
