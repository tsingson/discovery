#   configuration of Eureka / Discovery serverice daemon
discoveryd 即 Discovery Daemon 服务发现监听程序, fork 自 [bilibili/discovery](https://github.com/bilibili/discovery),  用于代码阅读学习与个人测试



该 fork 修改了
1. 可执行程序修改为 discoveryd , 源码在 /cmd/discoveryd/ 下, 并增加 daemon 自监听, 不需要 supersuvice 启动监听
2. 配置参数使用方式, 改为启动时读取指定的配置文件 discoveryd-config.toml
3. 配置参数中 scheduler 改为从 toml 配置文件中读取, 不再指定 scheduler.json 文件
4. scheduler.json 文件变更监听与热更换, 在将来实现 [TODO]



##  1. configuration file 配置文件
配置文件名为 discoveryd-config.toml
配置文件必须与 discoveryd 可执行文件位于同一路径下, 否则无法读配置文件, 导致 discoveryd 启动失败





discoveryd-config.toml
```
# 节点 IP , 可配置 discoveryd 其他节点, 形成集群, 多点相互同步 注册/发现 相关数据
Nodes = ["127.0.0.1:7171"]

# 当关 discoveryd 节点的配置, 用于自注册 appid = "infra.discovery" 的服务发现节点
[HTTPServer]
  Addr = "127.0.0.1:7171"

# 用于与其他 discoveryd 节点同步数据
[HTTPClient]
  Dial = "3s"
  KeepAlive = "120s"

# 最重要的配置项, 指定当前 discoveryd 节点所属区域 / 可用区 / host name 与 环境类型 , 与 HTTPServer.Addr 参数, 构成当前节点自注册的 5 个参数
# Env 项目配置不正确,  discovery 将无法正常提供服务(  clients 将匹配不到数据 )
[Env]
  Region = "china"
  Zone = "gd"
  Host = "discovery"
  DeployEnv = "dev"

# 初始化, 用于调度:  先匹配 Region 区域, 再匹配 Zone 可用区, 最后在 Env 运行环境中找 appid 与 hostname 指定的节点
# 有其他匹配的可选参数, 如 color ........ 
[Schedulers]
  [Schedulers."test.service-dev"]
    AppID = "test.service"
    Env = "dev"
    Remark = "test"

    [[Schedulers."test.service-dev".Zones]]
      Src = "gd"
      [Schedulers."test.service-dev".Zones.Dst]
        sz01 = 3

```



## 2. discoveryd 启动过程中的流程

discoveryd 即 Discovery Daemon 服务发现监听程序
