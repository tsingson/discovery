package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	json "github.com/json-iterator/go"

	"github.com/sanity-io/litter"

	"github.com/tsingson/discovery/conf"
)

func main() {
	cfg := &conf.Config{}
	cfg = conf.Default()
	data, err := json.Marshal(cfg)
	if err != nil {
		os.Exit(-1)
	}
	litter.Dump(string(data))

	fmt.Println("------------------->")
	cfg1 := &conf.Config{}

	err = json.Unmarshal(data, &cfg1)
	if err != nil {
		os.Exit(-1)
	}
	litter.Dump(cfg1)

	fmt.Println("------------------->")
	fmt.Println("------------------->")
	fmt.Println("------------------->")
	fmt.Println("------------------->")
	fmt.Println("------------------->")
	// data, err = yaml.Marshal(cfg)
	// if err != nil {
	// 	os.Exit(-1)
	// }
	// litter.Dump(string(data))
	// cfg2 := &conf.Config{}
	//
	// err = yaml.Unmarshal(data, &cfg2)
	// if err != nil {
	// 	os.Exit(-1)
	// }
	// litter.Dump(cfg2)

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(cfg); err != nil {
		os.Exit(-1)
	}
	fmt.Println(buf.String())

	cfg3 := &conf.Config{}
	b := buf.Bytes()
	err = toml.Unmarshal(b, &cfg3)
	if err != nil {
		os.Exit(-1)
	}

	litter.Dump(cfg3)
}
