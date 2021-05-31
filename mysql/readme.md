## MySQL 插件描述

插件使用xorm连接和操作数据库，详情见xorm.io。插件仅做了简单封装，旨在可以自动实例化数据库连接，并便捷使用。欢迎贡献想法
## MySQL 插件使用方式
1. 在plugin.yaml中写入如下配置
```yaml
clients:
  mysql:
      tars.test2.db:
        db_name: db_test
        user: root
        password: "123456"
        address: 127.0.0.1:3306
        args: charset=utf8mb4&loc=Asia%2FShanghai&allowOldPasswords=true&parseTime=true
      tars.test.db:
        db_name: db_test2
        user: root
        password: "123456"
        address: 127.0.0.1:3306
        args: charset=utf8mb4&loc=Asia%2FShanghai&allowOldPasswords=true&parseTime=true
```
2. 使用时需要预先在main.go中进行注册

```golang
package main

import (
    pluginManager "github.com/franklee0817/TarsGoPlugins/manager"
    _ "github.com/franklee0817/TarsGoPlugins/mysql"
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
import "github.com/franklee0817/TarsGoPlugins/mysql"

sess, err := mysql.NewSession(dbconf.TestDbInst)

```