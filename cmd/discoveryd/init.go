package main

import (
	"os"

	"github.com/sanity-io/litter"
	"github.com/spf13/afero"

	"github.com/tsingson/discovery/conf"
)

var ( // global variable

	path, logPath string
	cfg           *conf.Config
	configToml    string
)

func initx() {
	var err error
	afs := afero.NewOsFs()
	{ // setup path for storage of log / configuration / cache
		path = "/Users/qinshen/go/bin" // for test
		// path, err = file.GetCurrentExecDir()
		// if err != nil {
		// 	// fmt.Println("无法读取可执行程序的存储路径")
		// 	panic("无法读取可执行程序的存储路径")
		// 	os.Exit(-1)
		// }

	}
	{ // load config for discovery daemon
		configToml = path + "/discoveryd-config.toml"

		cfg = conf.Conf
		cfg, err = conf.LoadConfig(configToml)
		if err != nil {
			// fmt.Println("无法读取可执行程序的存储路径")
			//	panic("无法读取可执行程序的存储路径")
			os.Exit(-1)
		}
		litter.Dump(cfg)

	}
	{
		logPath = path + "/log"
		check, _ := afero.DirExists(afs, logPath)
		if !check {
			err = afs.MkdirAll(logPath, 0755)
			if err != nil {
				//	panic("mkdir log path fail")
				os.Exit(-1)
			}
		}
	}
}
