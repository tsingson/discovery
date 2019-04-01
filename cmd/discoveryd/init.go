package main

import (
	"os"

	"github.com/spf13/afero"

	"github.com/tsingson/discovery/conf"
	"github.com/tsingson/discovery/lib/file"
)

var ( // global variable


	// cacheSize                int
	// cacheTimeOut             int64

	path, logPath string
)

var cfg = conf.Conf

func init() {
	var err error
	afs := afero.NewOsFs()
	{ // setup path for storage of log / configuration / cache
		// path = "/Users/qinshen/git/linksmart/bin"  // for test
		path, err = file.GetCurrentExecDir()
		if err != nil {
			// fmt.Println("无法读取可执行程序的存储路径")
			panic("无法读取可执行程序的存储路径")
			os.Exit(-1)
		}

	}
	{ // load config for discovery daemon
		configToml := path + "/discoveryd-config.toml"

		cfg, err = conf.LoadConfig(configToml)
		if err != nil {
			// fmt.Println("无法读取可执行程序的存储路径")
			panic("无法读取可执行程序的存储路径")
			os.Exit(-1)
		}

	}
	{
		logPath = path + "/log"
		check, _ := afero.DirExists(afs, logPath)
		if !check {
			err = afs.MkdirAll(logPath, 0755)
			if err != nil {
				panic("mkdir log path fail")
				os.Exit(-1)
			}
		}
	}
}
