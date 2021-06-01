package redis

import (
	"errors"
	"sync"
	"time"

	"github.com/franklee0817/TarsGoPlugins/manager"
	redigo "github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v3"
)

const (
	PluginType = "clients"
	PluginName = "redis"
)

type clients struct {
	mux   *sync.RWMutex
	Pools map[string]*redigo.Pool
}

type config struct {
	MaxActive   int    `yaml:"max_active"`   //  最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
	MaxIdle     int    `yaml:"max_idle"`     // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
	IdleTimeout int    `yaml:"idle_timeout"` // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
	Wait        bool   `yaml:"wait"`         // 当超过最大连接数 是报错还是等待， true 等待 false 报错
	Address     string `yaml:"address"`      // redis 地址
	Db          int    `yaml:"Db"`           // redis数据库
}

var cfgs = make(map[string]config)
var Clients *clients

func init() {
	once := new(sync.Once)
	once.Do(func() {
		cfgs = make(map[string]config)
		if Clients == nil {
			Clients = new(clients)
			Clients.mux = new(sync.RWMutex)
			Clients.Pools = make(map[string]*redigo.Pool)
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
		pool := &redigo.Pool{
			MaxActive:   cfg.MaxActive,
			MaxIdle:     cfg.MaxIdle,
			IdleTimeout: time.Duration(cfg.IdleTimeout) * time.Second,
			Wait:        cfg.Wait,
			Dial: func() (redigo.Conn, error) {
				conn, err := redigo.Dial("tcp", cfg.Address, redigo.DialDatabase(cfg.Db))
				if err != nil {
					return nil, err
				}
				return conn, nil
			},
		}
		c.Pools[instName] = pool
	}

	return nil
}

// GetPool 获取指定实例名的连接池
func GetPool(instName string) (*redigo.Pool, error) {
	Clients.mux.RLock()
	defer Clients.mux.RUnlock()

	pool, ok := Clients.Pools[instName]
	if !ok {
		return nil, errors.New("no redis pool found for instance " + instName)
	}

	return pool, nil
}
