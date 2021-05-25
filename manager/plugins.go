package manager

import (
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/TarsCloud/TarsGo/tars"
)

// PluginCreator plugin的接口定义
type PluginCreator interface {
	Type() string
	Name() string
	Setup(cfg *yaml.Node) error
}

// PluginConfig 插件配置，type=>name=>config
type PluginConfig map[string]map[string]yaml.Node

var plugins = make(map[string]map[string]PluginCreator)
var pluginCfg PluginConfig

func init() {
	once := new(sync.Once)
	once.Do(func() {
		err := readCfg()
		if err != nil {
			panic(fmt.Errorf("init plugins failed:%v", err))
		}
	})
}

func readCfg() error {
	servCfg := tars.GetServerConfig()
	remoteConf := tars.NewRConf(servCfg.App, servCfg.Server, servCfg.BasePath)
	cfgStr, err := remoteConf.GetConfig("plugins.yaml")
	if err != nil {
		return err
	}

	var cfg = make(PluginConfig)
	err = yaml.Unmarshal([]byte(cfgStr), &cfg)
	if err != nil {
		return err
	}
	pluginCfg = cfg
	return nil
}

// Register 注册插件实例化方法
func Register(c PluginCreator) {
	if _, ok := plugins[c.Type()]; !ok {
		plugins[c.Type()] = make(map[string]PluginCreator)
	}
	plugins[c.Type()][c.Name()] = c
}

// GetPluginCreator 获取插件实例化方法
func GetPluginCreator(typ, name string) PluginCreator {
	if creators, ok := plugins[typ]; ok {
		if c, ok := creators[name]; ok {
			return c
		}
	}

	return nil
}

// Setup 初始化相关插件
func Setup() error {
	for typ, cfgs := range pluginCfg {
		for name, cfg := range cfgs {
			c := GetPluginCreator(typ, name)
			if c == nil {
				return fmt.Errorf("plugin creator has not been found for type:%v name:%v", typ, name)
			}
			err := c.Setup(&cfg)
			if err != nil {
				return fmt.Errorf("failed to setup plugin type:%v name:%v, err:%v", typ, name, err)
			}
		}
	}

	return nil
}