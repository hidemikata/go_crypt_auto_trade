package trade_manager

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"btcanallive_refact/config"
	"fmt"
	"reflect"
	"time"
)

type backtest_pc_table struct {
	test_num    int
	rci         int
	sma_long    int
	sma_short   int
	sma_up_rate int
}

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
func get_test_params() []backtest_pc_table {

	backtest_pc_number_use_table := make([]backtest_pc_table, 0)

	test_num := 1
	for rci_test_param := 5; rci_test_param <= 100; rci_test_param++ {
		for sma_long_i := 20; sma_long_i <= 50; sma_long_i++ {
			for sma_short_i := 3; sma_short_i <= 18; sma_short_i++ {
				for sma_up_rate := 4; sma_up_rate <= 7; sma_up_rate++ {
					param := backtest_pc_table{test_num, rci_test_param, sma_long_i, sma_short_i, sma_up_rate}
					backtest_pc_number_use_table = append(backtest_pc_number_use_table, param)
					test_num++
				}
			}
		}
	}

	my_test_params := make([]backtest_pc_table, 0)
	for i, v := range backtest_pc_number_use_table {
		tme_devide := backtest_pc_number_use_table[i].test_num % config.Config.NumOfPc
		if tme_devide == 0 && config.Config.NumOfPc == config.Config.PcNoumber {
			my_test_params = append(my_test_params, v)
		} else if tme_devide == config.Config.PcNoumber {
			my_test_params = append(my_test_params, v)
		} else {

		}
	}
	return my_test_params
}
func BacktestStart(ti []trade_jadge_algo.TradeInterface) {

	test_parmas := get_test_params()

	fmt.Println("backteststart")

	loc, _ := time.LoadLocation("Asia/Tokyo")
	alna_minute_max := 23040
	//alna_minute_max = 440

	for _, param_v := range test_parmas {
		set_test_param("*trade_jadge_algo.Rci", ti, param_v.rci)
		fmt.Println("start rci = :", param_v.rci, "time=", time.Now())

		set_test_param("*trade_jadge_algo.Sma", ti, param_v.sma_long, param_v.sma_short, param_v.sma_up_rate)
		model.ClearBacktestPosition()
		var profit float64
		var position_count int
		profit = 0.0
		position_count = 0
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
		profit, position_count = model.GetProfitBacktest()
		fmt.Println("profit=", profit)

		model.BacktestInsertTotalProfit(time.Now(), profit, param_v.sma_long, param_v.sma_short, float64(param_v.sma_up_rate)/1000, param_v.rci, position_count)
	}
}

var layout = "2006-01-02 15:04:05"

func timeToString(t time.Time) string {
	str := t.Format(layout)
	return str
}
