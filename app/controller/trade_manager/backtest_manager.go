package trade_manager

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"btcanallive_refact/config"
	"fmt"
	"time"
)

func is_my_test(test_num int) bool {

	test_num = test_num + 1
	target_num := test_num % config.Config.NumOfPc
	if target_num == 0 && config.Config.NumOfPc == config.Config.PcNoumber {
		return true
	} else if target_num == config.Config.PcNoumber {
		return true
	} else {
		return false
	}
}
func create_all_set(all [][]trade_jadge_algo.BacktestParams, p []trade_jadge_algo.BacktestParams) [][]trade_jadge_algo.BacktestParams {
	all_params := make([][]trade_jadge_algo.BacktestParams, 0)
	for _, v_all := range all {
		for _, v_p := range p {
			tmp := v_all
			tmp = append(tmp, v_p)
			all_params = append(all_params, tmp)
		}
	}
	return all_params
}
func get_test_params(ti []trade_jadge_algo.TradeInterface) [][]trade_jadge_algo.BacktestParams {

	ti_params := make([][]trade_jadge_algo.BacktestParams, 0)
	all_params := make([][]trade_jadge_algo.BacktestParams, 0)

	for _, ti_v := range ti {
		ti_params = append(ti_params, ti_v.CreateBacktestParams())
	}

	for _, vv := range ti_params[0] {
		all_params = append(all_params, []trade_jadge_algo.BacktestParams{vv})
	}
	if len(ti_params) >= 2 {
		for i := 1; i < len(ti_params); i++ {
			all_params = create_all_set(all_params, ti_params[i])
		}
	}

	return all_params
}
func BacktestStart(ti []trade_jadge_algo.TradeInterface) {
	fmt.Println("backteststart")

	test_paramas := get_test_params(ti)

	loc, _ := time.LoadLocation("Asia/Tokyo")
	alna_minute_max := 1000

	for test_i, test_v := range test_paramas {
		fmt.Println(test_i, test_v)
	}

	for param_i, param_v := range test_paramas {
		if !is_my_test(param_i) {
			continue
		}

		fmt.Println("test count", param_i+1, "/", len(test_paramas))

		for _, p_v := range param_v {
			p_v.Ti.BacktestSetParam(p_v.Params)
		}

		model.ClearBacktestPosition()
		var profit float64
		var position_count int
		profit = 0.0
		position_count = 0
		for time_i := 0; time_i < alna_minute_max; time_i++ {
			anal_time := time.Date(2021, 7, 29, 6, time_i, 0, 0, loc)
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
			}

			latest_pos, fixed := model.GetLatestPosition(true)

			candle_date := model.GetLatestCandle(timeToString(anal_time))

			tick := trade_def.Ticker{}
			tick.BestAsk = candle_date.Close
			tick.BestBid = candle_date.Close - float64(config.Config.Spread) //spread

			if buy && (latest_pos.Date == "" || fixed) {
				fmt.Println("test buy")
				model.InsertPosition(anal_time, tick, true)
			} else if fix && (latest_pos.Date != "" && !fixed) {
				fmt.Println("test fix")
				model.UpdatePosition(anal_time, latest_pos, tick, true)
			} else {
			}
		}
		profit, position_count = model.GetProfitBacktest()
		fmt.Println("profit=", profit, "testcount", param_i)
		position_count = 0
		//コンパイルエラーなので殺してる
		//		model.BacktestInsertTotalProfit(time.Now(), profit, param_v.sma_long, param_v.sma_short, float64(param_v.sma_up_rate)/1000, param_v.rci, position_count, param_v.rci_test_buy_rate)
		model.BacktestInsertTotalProfit(time.Now(), profit, 0, 0, float64(0)/1000, 0, position_count, 0)
	}
}

var layout = "2006-01-02 15:04:05"

func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}
