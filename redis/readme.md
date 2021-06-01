## redis 插件描述

插件使用`github.com/gomodule/redigo/redis`连接和操作redis。插件仅做了简单封装，旨在可以自动实例化redis连接，并便捷使用。欢迎贡献想法
## redis 插件使用方式
1. 在plugin.yaml中写入如下配置
```yaml
clients:
  redis:
      tars.redis.test1:
        max_active: 100
        max_idle: 100
        idle_timeout: 100
        wait: true
        address: "127.0.0.1:6379"
        Db: 1
      tars.redis.test2:
        max_active: 100
        max_idle: 100
        idle_timeout: 100
        wait: true
        address: "127.0.0.1:6379"
        Db: 2
```
2. 使用时需要预先在main.go中进行注册

```golang
package main

import (
    pluginManager "github.com/franklee0817/TarsGoPlugins/manager"
    _ "github.com/franklee0817/TarsGoPlugins/redis"
)

func main() {
	...
	err := pluginManager.Setup()
	if err != nil {
		panic(err)
    }
	...
}
```

3. 获取数据库连接会话
```golang
import "github.com/franklee0817/TarsGoPlugins/redis"

pool, err := redis.GetPool(dbconf.TestRedisInst)

```