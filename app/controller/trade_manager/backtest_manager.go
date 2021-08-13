package trade_manager

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"fmt"
	"reflect"
	"time"
)

func set_test_param(ti_string string, ti []trade_jadge_algo.TradeInterface, param int) {
	for _, ti_v := range ti {
		if reflect.TypeOf(ti_v).String() == ti_string {
			ti_v.SetParam(param)
			return
		}
	}
	panic("not found trade interface")
}

func BacktestStart(ti []trade_jadge_algo.TradeInterface) {
	fmt.Println("backteststart")

	loc, _ := time.LoadLocation("Asia/Tokyo")
	alna_minute_max := 23040

	//rciloop
	for rci_test_param := 5; rci_test_param <= 100; rci_test_param++ {
		set_test_param("*trade_jadge_algo.Rci", ti, rci_test_param)
		//ここでsma loopする予定。
		fmt.Println("start rci = :", rci_test_param, "time=", time.Now())

		model.ClearBacktestPosition()
		for time_i := 0; time_i < alna_minute_max; time_i++ {
			anal_time := time.Date(2021, 7, 28, 6, time_i, 0, 0, loc)
			//fmt.Println("start:", anal_time)

			buy := true
			fix := false
			for _, ti_v := range ti {
				if !ti_v.IsDbCollectedData(anal_time) {
					buy = false
					fix = false
					continue
				}

				ti_v.Analisis(anal_time)
				buy = buy && ti_v.IsTradeOrder()
				fix = fix || ti_v.IsTradeFix()
				//fmt.Println("test", buy, fix)
			}

			latest_pos, fixed := model.GetLatestPosition(true)

			candle_date := model.GetLatestCandle(timeToString(anal_time))

			tick := trade_def.Ticker{}
			tick.BestAsk = candle_date.Close
			tick.BestBid = candle_date.Close - 2000 //spread

			if buy && (latest_pos.Date == "" || fixed) {
				fmt.Println("test buy")
				model.InsertPosition(anal_time, tick, true)
			} else if fix && (latest_pos.Date != "" && !fixed) {
				fmt.Println("test fix")
				model.UpdatePosition(anal_time, latest_pos, tick, true)
			} else {
			}
		}
		profit := model.GetProfitBacktest()
		fmt.Println("profit=", profit)

		model.BacktestInsertTotalProfit(time.Now(), profit, 25, 5, 0.005, rci_test_param)
	}
}

var layout = "2006-01-02 15:04:05"

func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}
