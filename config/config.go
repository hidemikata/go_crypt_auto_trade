package config

import (
	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ApiKey             string
	ApiSecret          string
	BackTest           bool
	NumOfPc            int
	PcNoumber          int
	BackTestInMemory   bool
	Spread             int
	SmaLong            int
	SmaLongLong        int
	SmaShort           int
	SmaUpToRate        float64
	SmaUpToRateLatest1 float64
	RciLong            int
	RciRate            float64
}

var Config ConfigList

func init() {
	cfg, _ := ini.Load("config/config.ini")
	Config = ConfigList{
		ApiKey:           cfg.Section("bitflyer").Key("api_key").String(),
		ApiSecret:        cfg.Section("bitflyer").Key("secret_key").String(),
		BackTest:         cfg.Section("backtest").Key("backtest").String() == "true",
		NumOfPc:          cfg.Section("backtest").Key("num_of_pc").MustInt(),
		PcNoumber:        cfg.Section("backtest").Key("pc_number").MustInt(),
		BackTestInMemory: cfg.Section("backtest").Key("backtest_inmemory").String() == "true" && cfg.Section("backtest").Key("backtest").String() == "true",
		Spread:           cfg.Section("backtest").Key("spread").MustInt(),

		SmaLong:            cfg.Section("analisys").Key("sma_long").MustInt(),
		SmaLongLong:        cfg.Section("analisys").Key("sma_long_long").MustInt(),
		SmaShort:           cfg.Section("analisys").Key("sma_short").MustInt(),
		SmaUpToRate:        cfg.Section("analisys").Key("sma_up_to_rate").MustFloat64(),
		SmaUpToRateLatest1: cfg.Section("analisys").Key("sma_up_to_rate_latest1").MustFloat64(),
		RciLong:            cfg.Section("analisys").Key("rci_long").MustInt(),
		RciRate:            cfg.Section("analisys").Key("rci_rate").MustFloat64(),
	}
}
