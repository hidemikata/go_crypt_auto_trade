package trade_manager

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"fmt"
	"reflect"
	"time"
)

func set_test_param(ti_string string, ti []trade_jadge_algo.TradeInterface, param ...int) {

	for _, ti_v := range ti {
		if reflect.TypeOf(ti_v).String() == "*trade_jadge_algo.Rci" && ti_string == "*trade_jadge_algo.Rci" {
			ti_v.SetParam(param[0])
			return
		} else if reflect.TypeOf(ti_v).String() == "*trade_jadge_algo.Sma" && ti_string == "*trade_jadge_algo.Sma" {
			ti_v.SetParam(param[0], param[1], param[2])
			return
		}
	}
	panic("not found trade interface")
}

func BacktestStart(ti []trade_jadge_algo.TradeInterface) {
	fmt.Println("backteststart")

	loc, _ := time.LoadLocation("Asia/Tokyo")
	alna_minute_max := 23040
	//alna_minute_max = 440

	//rciloop
	for rci_test_param := 51; rci_test_param <= 51; rci_test_param++ {
		set_test_param("*trade_jadge_algo.Rci", ti, rci_test_param)
		fmt.Println("start rci = :", rci_test_param, "time=", time.Now())

		for sma_long_i := 26; sma_long_i <= 26; sma_long_i++ {
			for sma_short_i := 10; sma_short_i <= 10; sma_short_i++ {
				for sma_up_rate := 5; sma_up_rate <= 5; sma_up_rate++ {
					set_test_param("*trade_jadge_algo.Sma", ti, sma_long_i, sma_short_i, sma_up_rate)
					model.ClearBacktestPosition()
					for time_i := 0; time_i < alna_minute_max; time_i++ {
						anal_time := time.Date(2021, 7, 28, 6, time_i, 0, 0, loc)
						//anal_time = time.Date(2021, 8, 2, 23, time_i, 0, 0, loc)
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
					profit, position_count := model.GetProfitBacktest()
					fmt.Println("profit=", profit)

					model.BacktestInsertTotalProfit(time.Now(), profit, sma_long_i, sma_short_i, float64(sma_up_rate)/1000, rci_test_param, position_count)
				}
			}
		}
	}
}

var layout = "2006-01-02 15:04:05"

func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}
