package trade_manager

import (
	"btcanallive_refact/app/bitflyer"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/marcket_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"btcanallive_refact/app/model"
	"fmt"
	"time"
)

func Run() {
	fmt.Println("run")
    api := bitflyer.NewBitflyer("BTC_JPY")
    var marcket marcket_def.Marcket
    marcket = api
	real_time_ticker_ch := make(chan trade_def.Ticker, 1)
	go marcket.StartGettingRealTimeTicker(real_time_ticker_ch)
	defer close(real_time_ticker_ch)

	for i := range real_time_ticker_ch {
		if !save_ticker_table(i) {
			continue
		}

		sma_algo := trade_jadge_algo.NewSmaAlgorithm()
		var ti trade_jadge_algo.TradeInterface
		ti = sma_algo
		if !ti.IsDbCollectedData() {
			continue
		}
        ti.Analisis()
		buy := ti.IsTradeOrder()
		fix := ti.IsTradeFix()
        fmt.Println(buy, fix)
        latest_pos, fixed:= model.GetLatestPosition()
		if buy && (latest_pos.Date == "" || fixed) {
            fmt.Println("buy")
			marcket.PutOrder()
            tick := marcket.GetTicker()
            model.InsertPosition(tick)
		} else if fix && (latest_pos.Date != "" && !fixed){
            fmt.Println("fix")
			marcket.FixOrder()
            tick := marcket.GetTicker()
            model.UpdatePosition(latest_pos, tick)
        } else {
        }
	}

}

func second_to_zero(t time.Time) string {
	min := fmt.Sprintf("%02d", t.Minute())
	h := fmt.Sprintf("%02d", t.Hour())
	d := fmt.Sprintf("%02d", t.Day())
	m := fmt.Sprintf("%02d", int(t.Month()))
	y := fmt.Sprintf("%02d", t.Year())
	return y + "-" + m + "-" + d + " " + h + ":" + min + ":00"
}

func truncate_minute(date string) string {
	t := timeComvertAdd9hour(date)
	return second_to_zero(t)
}

func timeComvertAdd9hour(date string) time.Time {
	t, _ := time.Parse(time.RFC3339, date)
	return t.Add(time.Hour * 9)
}

func save_ticker_table(t trade_def.Ticker) bool {
	insert := false
	fmt.Print(".")
	date := truncate_minute(t.Timestamp)

	count := model.GetNumOfCandle(date)

	if count == 0 {
		fmt.Println("insert")
		model.InsertNewCandle(date, t)
		insert = true
	} else {
		bj := model.GetLatestCandle(date)
		h := bj.High
		l := bj.Low
		c := t.BestAsk
		if bj.High < t.BestAsk {
			h = t.BestAsk
		} else if bj.Low > t.BestAsk {
			l = t.BestAsk
		}
		model.UpdateCandle(date, h, l, c)
	}
	return insert
}

