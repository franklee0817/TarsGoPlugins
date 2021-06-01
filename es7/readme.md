## es 插件描述

插件使用`github.com/olivere/elastic`连接和操作es。插件仅做了简单封装，旨在可以自动实例化ES连接，并便捷使用。欢迎贡献想法
## es 插件使用方式
1. 在plugin.yaml中写入如下配置
```yaml
clients:
  es7:
      tars.es.test:
        url: "http://127.0.0.1:9200"
        user: "root"
        password: "123456"
        sniff: false
        health_check: false
      tars.es.test2:
        url: "http://127.0.0.1:9200"
        user: "root"
        password: "123456"
```
2. 使用时需要预先在main.go中进行注册

```golang
package main

import (
    pluginManager "github.com/franklee0817/TarsGoPlugins/manager"
    _ "github.com/franklee0817/TarsGoPlugins/es7"
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
import "github.com/franklee0817/TarsGoPlugins/es7"

client, err := es.GetClient(esconf.TestESInst)

```