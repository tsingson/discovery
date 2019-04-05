# note 说明

fork from [https://github.com/bilibili/discovery](https://github.com/bilibili/discovery)

这是一个用于学习的 fork, 用于与 goim 的配合测试,  代码修改了以下:

 * [x] 去除所有命令行参数, 改为读取 toml 配置文件
 * [x] 可执行程序改为 daemon 后台运行, 代码位于 /cmd/discoveryd/ 下
 * [ ] 修改注册数据为可从文件读取的缓存, 以便 discovery 重启时自动恢复原运行数据



***代码随时变更, 无法保障使用***






###  discoveryd 运行

请下载 [编译好的discoveryd用于测试](https://github.com/tsingson/discovery/releases/download/v0.1.1/discoveryd.tar.bz2)  并用以下命令解码到 linux ( 验证环境为 cent os 7 ) 到某一路径下 

	注: 压缩包的打包方式 tar -cvjpf discoveryd.tar.bz2 ./discoveryd


```
tar -xvjf    discoveryd.tar.bz2
cd discoveryd
chmod +x ./discoveryd
ps -ef | grep discoveryd
```

 运行, 注: 该程序为后台 daemon 运行, 用 ps -ef | grep discoveryd 查看是否成功运行

```
discoveryd
```



----------
# Discovery 

[![Build Status](https://travis-ci.org/bilibili/discovery.svg?branch=master)](https://travis-ci.org/bilibili/discovery) 
[![Go Report Card](https://goreportcard.com/badge/github.com/bilibili/discovery)](https://goreportcard.com/report/github.com/bilibili/discovery)
[![codecov](https://codecov.io/gh/Bilibili/discovery/branch/master/graph/badge.svg)](https://codecov.io/gh/Bilibili/discovery)

Discovery is a based service that is production-ready and primarily used at [Bilibili](https://www.bilibili.com/) for locating services for the purpose of load balancing and failover of middle-tier servers.

## Quick Start

### env

`go1.9.x` (and later)

### build
```shell
cd $GOPATH/src
git clone https://github.com/bilibili/discovery.git
cd discovery/cmd/discovery
go build
```

### run
```shell
./discovery -conf discovery-example.toml -alsologtostderr
```

`-alsologtostderr` is `glog`'s flag，means print into stderr. If you hope print into file, can use `-log_dir="/tmp"`. [view glog doc](https://godoc.org/github.com/golang/glog).

### Configuration

You can view the comments in `cmd/discovery/discovery-example.toml` to understand the meaning of the config.

### Client

* [API Doc](doc/api.md)
* [Go SDK](naming/client.go) | [Example](naming/example_test.go)
* [Java SDK](https://github.com/flygit/discoveryJavaSDK)
* [CPP SDK](https://github.com/brpc/brpc/blob/master/src/brpc/policy/discovery_naming_service.cpp)
* [other languaue](doc/sdk.md)

## Intro/Arch/Practice

* [Introduction](doc/intro.md)
* [Architecture](doc/arch.md)
* [Practice in Bilibili](doc/practice.md)

## Feedback

Please report bugs, concerns, suggestions by issues, or join QQ-group 716486124 to discuss problems around source code.
