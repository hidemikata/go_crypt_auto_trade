package trade_manager

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"fmt"
	"time"
)

func BacktestStart(ti []trade_jadge_algo.TradeInterface) {
	fmt.Println("backteststart")

	model.ClearBacktestPosition()

	loc, _ := time.LoadLocation("Asia/Tokyo")
	alna_minute_max := 129600 //90日
	alna_minute_max = 1000
	for time_i := 0; time_i < alna_minute_max; time_i++ {
		anal_time := time.Date(2021, 8, 10, 8, time_i, 0, 0, loc)
		fmt.Println(anal_time)

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

		ask, bid := ti[0].GetLatestAskBid() //暫定でsmaから取る
		tick := trade_def.Ticker{}
		tick.BestAsk = ask
		tick.BestBid = bid

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
}
