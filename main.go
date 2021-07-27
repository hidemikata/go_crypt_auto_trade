package main

import (
	"btcanallive_refact/app/bitflyer"
	"btcanallive_refact/app/controller/trade_manager"
	"btcanallive_refact/app/marcket_def"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"sync"
    "time"
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
    time_ch := make(chan bool, 1)
    defer close(time_ch)

    go time_checker(time_ch)

	var wg sync.WaitGroup
	wg.Add(2)
	go trade_manager.StartRealTimeTickGetter(marcket, real_time_ticker_ch)//こっちの通信も止めたいけど止めれない
	go trade_manager.StartAnalisis(marcket, real_time_ticker_ch, ti, time_ch)
	wg.Wait()
}

func time_checker(ch chan bool){
    for {
        time.Sleep(50 * time.Millisecond)
        if len(ch) != 0 {
            continue
        }
        not_time := time.Now()
	    h:=not_time.Hour()
        if 3 < h && h < 7{
            ch<-false
        }
        ch<-true
    }
}


