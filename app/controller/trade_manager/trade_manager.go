package trade_manager

import (
	"btcanallive_refact/app/marcket_def"
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"fmt"
	"sync"
	"time"
)

func StartRealTimeTickGetter(marcket marcket_def.Marcket, real_time_ticker_ch chan trade_def.Ticker, wg *sync.WaitGroup) {
	fmt.Println("StartRealTimeTickGetter")
	marcket.StartGettingRealTimeTicker(real_time_ticker_ch)
	wg.Done()
}

func StartAnalisis(marcket marcket_def.Marcket, real_time_ticker_ch chan trade_def.Ticker, ti []trade_jadge_algo.TradeInterface, time_ch chan bool) {
	fmt.Println("StartAnalisis")
	for i := range real_time_ticker_ch {
		if <-time_ch == false {
			ForceMarcketClose(marcket)
			continue
		}
		if !save_ticker_table(i) {
			continue
		}

		buy := true
		fix := false
		fmt.Println("0", buy, fix)
		for _, ti_v := range ti {
			if !ti_v.IsDbCollectedData() {
				buy = false
				fix = false
				continue
			}

			ti_v.Analisis()
			buy = buy && ti_v.IsTradeOrder()
			fix = fix || ti_v.IsTradeFix()
			fmt.Println("-", buy, fix)
		}

		fmt.Println("1", buy, fix)
		latest_pos, fixed := model.GetLatestPosition()

		if buy && (latest_pos.Date == "" || fixed) {
			fmt.Println("buy")
			marcket.PutOrder()
			tick := marcket.GetTicker()
			model.InsertPosition(tick)
		} else if fix && (latest_pos.Date != "" && !fixed) {
			fmt.Println("fix")
			marcket.FixOrder()
			tick := marcket.GetTicker()
			model.UpdatePosition(latest_pos, tick)
		} else {
		}
	}

}

func ForceMarcketClose(marcket marcket_def.Marcket) {
	fmt.Print("x")
	latest_pos, fixed := model.GetLatestPosition()
	if latest_pos.Date != "" && !fixed {
		fmt.Println("foce fix do")
		marcket.FixOrder()
		tick := marcket.GetTicker()
		model.UpdatePosition(latest_pos, tick)
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
		fmt.Print("\ninsert-")
		fmt.Print(time.Now(), "-")
		_, _, _, time_count := model.GetCandleData()
		fmt.Println(time_count)
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
