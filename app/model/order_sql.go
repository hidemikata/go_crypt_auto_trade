package model

import (
	"btcanallive_refact/app/trade_def"
	"fmt"
	"strconv"
	"time"
)

func InsertPosition(now time.Time, tick trade_def.Ticker, is_backtest bool) {

	date_str := now.Format(format1)
	table_name := "btc_jpy_live_position"
	if is_backtest {
		table_name = "btc_jpy_live_position_backtest"
	}
	query := `insert into ` + table_name + ` values("` + date_str + `", "` + "buy" +
		`", ` + strconv.FormatFloat(tick.BestAsk, 'f', -1, 64) + `, NULL` +
		`, NULL` +
		`, NULL` +
		`, "` + "BTC_JPY" +
		`");`

	fmt.Println("insert_positon_query:", query)
	_, err := db.Exec(query)

	if err != nil {
		panic(err.Error())
	}
}

func UpdatePosition(now time.Time, pos trade_def.Position, tick trade_def.Ticker, is_backtest bool) {
	fix_date := now.Format(format1)

	pos.Fix_date = fix_date
	pos.Fix_price = tick.BestBid
	pos.Profit = tick.BestBid - pos.Price

	table_name := "btc_jpy_live_position"
	if is_backtest {
		table_name = "btc_jpy_live_position_backtest"
	}
	query := `update ` + table_name + ` set fix_date="` + pos.Fix_date +
		`", fix_price=` + strconv.FormatFloat(pos.Fix_price, 'f', -1, 64) +
		`, profit=` + strconv.FormatFloat(pos.Profit, 'f', -1, 64) +
		` where date="` + pos.Date + `";`

	_, err := db.Exec(query)

	fmt.Println("fix:", query)

	if err != nil {
		panic(err.Error())
	}

}

func ClearBacktestPosition() {
	query := "delete from btc_jpy_live_position_backtest;"

	_, err := db.Exec(query)

	fmt.Println("delete:", query)

	if err != nil {
		panic(err.Error())
	}
}
