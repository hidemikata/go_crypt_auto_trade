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
	//$B%"%k%4$r9%$-$J$@$1(Bnew$B$7$F(Bappend$B$9$k(B
	sma_algo := trade_jadge_algo.NewSmaAlgorithm()
	ti := make([]trade_jadge_algo.TradeInterface, 0)
	ti = append(ti, sma_algo)


    //$B;~4V$,Mh$?$i$V$C$?@Z$C$?$j:F3+$7$?$j!#(B

	var wg sync.WaitGroup
	wg.Add(2)
	go trade_manager.StartRealTimeTickGetter(marcket, real_time_ticker_ch)
	go trade_manager.StartAnalisis(marcket, real_time_ticker_ch, ti)
	wg.Wait()
}
