package mysql

import (
	"sync"

	_ "github.com/go-sql-driver/mysql"

	"github.com/franklee0817/TarsGoPlugins/manager"

	"gopkg.in/yaml.v3"
	"xorm.io/xorm"
)

const (
	PluginType = "clients"
	PluginName = "mysql"
)

// client mysql client 结构
type client struct {
	Engine *xorm.Engine
}

var Client *client

func init() {
	once := new(sync.Once)
	once.Do(func(){
		if Client == nil {
			Client = new(client)
		}
		manager.Register(PluginType, Client)
	})
}

// Type 获取插件类型
func (m *client) Type() string {
	return PluginType
}

// Name 获取插件名称
func (m *client) Name() string {
	return PluginName
}

// Setup 初始化插件
func (m *client) Setup(cfg *yaml.Node) error {
	return nil
}
