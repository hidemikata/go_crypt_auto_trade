package config

import (
	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey           string
	ApiSecret        string
	BackTest         string
	NumOfPc          int
	PcNoumber        int
	BackTestInMemory string
}

var Config ConfigList

func init() {
	cfg, _ := ini.Load("config/config.ini")
	Config = ConfigList{
		ApiKey:           cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret:        cfg.Section("bitflyer").Key("secret_key").String(),
		BackTest:         cfg.Section("backtest").Key("backtest").String(),
		NumOfPc:          cfg.Section("backtest").Key("num_of_pc").MustInt(),
		PcNoumber:        cfg.Section("backtest").Key("pc_number").MustInt(),
		BackTestInMemory: cfg.Section("backtest").Key("backtest_inmemory").String(),
	}
}
