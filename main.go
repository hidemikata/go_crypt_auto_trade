package main

import (
	"btcanallive_refact/app/bitflyer"
	"btcanallive_refact/app/controller/trade_manager"
	"btcanallive_refact/app/marcket_def"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"sync"
)

func main() {

	api := bitflyer.NewBitflyer("BTC_JPY")
	var marcket marcket_def.Marcket
	marcket = api

	real_time_ticker_ch := make(chan trade_def.Ticker, 1)
	defer close(real_time_ticker_ch)
	//アルゴを好きなだけnewしてappendする
	sma_algo := trade_jadge_algo.NewSmaAlgorithm()
	ti := make([]trade_jadge_algo.TradeInterface, 0)
	ti = append(ti, sma_algo)


    //時間が来たらぶった切ったり再開したり。

	var wg sync.WaitGroup
	wg.Add(2)
	go trade_manager.StartRealTimeTickGetter(marcket, real_time_ticker_ch)
	go trade_manager.StartAnalisis(marcket, real_time_ticker_ch, ti)
	wg.Wait()
}
