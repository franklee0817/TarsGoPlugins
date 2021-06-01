package mysql

import (
	"errors"
	"fmt"
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

// clients mysql client 结构
type clients struct {
	mux     *sync.RWMutex
	Engines map[string]*xorm.Engine
}

type config struct {
	DBName   string `yaml:"db_name"`  // 数据库名称
	User     string `yaml:"user"`     // 用户名称
	Password string `yaml:"password"` // 用户密码
	Address  string `yaml:"address"`  // 数据库连接地址
	Args     string `yaml:"args"`     // 连接参数
}

// Clients 全局mysql连接的clients实例
var Clients *clients

var cfgs map[string]config

func init() {
	once := new(sync.Once)
	once.Do(func() {
		cfgs = make(map[string]config)
		if Clients == nil {
			Clients = new(clients)
			Clients.mux = new(sync.RWMutex)
			Clients.Engines = make(map[string]*xorm.Engine)
		}
		manager.Register(Clients)
	})
}

// Type 获取插件类型
func (m *clients) Type() string {
	return PluginType
}

// Name 获取插件名称
func (m *clients) Name() string {
	return PluginName
}

// Setup 初始化插件
func (m *clients) Setup(yamlCfg *yaml.Node) error {
	err := yamlCfg.Decode(&cfgs)
	if err != nil {
		return err
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	for instName, cfg := range cfgs {
		m.Engines[instName], err = xorm.NewEngine("mysql", buildDSN(cfg))
		if err != nil {
			return err
		}
	}
	return nil
}

func buildDSN(cfg config) string {
	dsn := fmt.Sprintf("%s:%s@(%s)/%s", cfg.User, cfg.Password, cfg.Address, cfg.DBName)
	if len(cfg.Args) > 0 {
		dsn += "?" + cfg.Args
	}

	return dsn
}

// NewSession 创建一个新的MySQL 连接会话
func NewSession(instName string) (*xorm.Session, error) {
	Clients.mux.RLock()
	engine := Clients.Engines[instName]
	if engine == nil {
		return nil, errors.New("no mysql instance found for " + instName)
	}
	err := engine.Ping()
	if err != nil {
		engine, err := refreshEngine(instName)
		if err != nil {
			return nil, errors.New("cannot connect to mysql instance " + instName)
		}
		Clients.mux.RUnlock()
		Clients.mux.Lock()
		Clients.Engines[instName] = engine
		return engine.NewSession(), nil
	}
	Clients.mux.RUnlock()

	return engine.NewSession(), nil
}

func refreshEngine(instName string) (*xorm.Engine, error) {
	cfg := cfgs[instName]
	connStr := fmt.Sprintf("%s:%s@(%s)/%s", cfg.User, cfg.Password, cfg.Address, cfg.DBName)
	if len(cfg.Args) > 0 {
		connStr += "?" + cfg.Args
	}
	return xorm.NewEngine("mysql", connStr)
}
