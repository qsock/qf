package main

import (
	"github.com/BurntSushi/toml"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/util/conf"
	"os"
)

func main() {
	path := "config.toml"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	if err := conf.LoadFromReader(f, toml.Unmarshal); err != nil {
		panic(err)
	}
	logCfg := qlog.Config{}
	if err := conf.UnmarshalKey("log", &logCfg); err != nil {
		panic(err)
	}
	logCfg.CallerSkip = 1
	qlog.SetCfg(&logCfg)
	qlog.Infof("val:%#v||normal:%v", logCfg)
}
