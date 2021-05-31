## Plugin Manager

这是一个插件管理包，通过这个包在启动时自动读取tconf的`plugin.yaml`配置，并注册相关插件，
通过显示的调用 manager.Setup()，实现所有插件的实例化。

暂未计划引入别的feature，如果有需要的话，欢迎提issue