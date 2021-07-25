package model

import (
	"btcanallive_refact/app/trade_def"
	"fmt"
	"strconv"
	"time"
)

func InsertPosition(tick trade_def.Ticker) {
	now := time.Now()
	date_str := now.Format(format1)

	query := `insert into btc_jpy_live_position values("` + date_str + `", "` + "buy" +
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

func UpdatePosition(pos trade_def.Position, tick trade_def.Ticker) {
	now := time.Now()
	fix_date := now.Format(format1)

	pos.Fix_date = fix_date
	pos.Fix_price = tick.BestBid
	pos.Profit = tick.BestBid - pos.Price

	query := `update btc_jpy_live_position set fix_date="` + pos.Fix_date +
		`", fix_price=` + strconv.FormatFloat(pos.Fix_price, 'f', -1, 64) +
		`, profit=` + strconv.FormatFloat(pos.Profit, 'f', -1, 64) +
		` where date="` + pos.Date + `";`

	_, err := db.Exec(query)

	fmt.Println("fix:", query)

	if err != nil {
		panic(err.Error())
	}

}
