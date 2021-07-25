package config

import (
    "gopkg.in/ini.v1"
)

type ConfigList struct {
    ApiKey     string
    ApiSecret string
}

var Config ConfigList

func init() {
    cfg, _ := ini.Load("config/config.ini")
    Config = ConfigList{
        ApiKey:     cfg.Section("bitflyer").Key("api_key").String(),
        ApiSecret:  cfg.Section("bitflyer").Key("secret_key").String(),
    }
}
