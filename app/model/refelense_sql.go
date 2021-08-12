package model

import (
	"btcanallive_refact/app/trade_def"
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

func GetLatestPosition(is_backtest bool) (trade_def.Position, bool) {

	table_name := "btc_jpy_live_position"
	if is_backtest {
		table_name = "btc_jpy_live_position_backtest"
	}

	rows, err := db.Query(`select * from ` + table_name + ` order by date desc limit 1;`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var pos trade_def.Position
	null := new(sql.NullString)

	for rows.Next() {
		err = rows.Scan(
			&pos.Date,
			&pos.Buy_or_sell,
			&pos.Price,
			null,
			null,
			null,
			&pos.Symbol,
		)
		if err != nil {
			panic(err.Error())
		}
		break
	}
	//null.Valid true = not null
	return pos, null.Valid
}

func GetNumberOfCandleBetweenDate(before_data_str string, now_str string) int {
	var count int
	err := db.QueryRow(`select count(*) from btc_jpy_live where date between "` + before_data_str + `" and "` + now_str + `";`).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	return count
}

func GetCandleBetweenDate(past_str string, latest_str string) []trade_def.BtcJpy {
	//データをDBから取得
	rows, err := db.Query(`select * from btc_jpy_live where date between "` + past_str + `" and "` + latest_str + `" order by date;`)
	if err != nil {
		panic(err.Error())
	}

	records := make([]trade_def.BtcJpy, 0)
	for rows.Next() {
		var record trade_def.BtcJpy
		err = rows.Scan(
			&record.Date,
			&record.Symbol,
			&record.Open,
			&record.High,
			&record.Low,
			&record.Close,
		)
		if err != nil {
			panic(err.Error())
		}
		records = append(records, record)
	}

	return records
}

func GetNumOfCandle(date string) int {

	var count int
	err := db.QueryRow(`select count(*) from btc_jpy_live where date="` + date + `";`).Scan(&count)
	if err != nil {
		panic(err.Error())
	}
	return count
}

func InsertNewCandle(date string, t trade_def.Ticker) {

	query_insert := `insert into btc_jpy_live values("` + date + `"` +
		`, "BTC_JPY"` +
		", " + strconv.FormatFloat(t.BestAsk, 'f', -1, 64) +
		", " + strconv.FormatFloat(t.BestAsk, 'f', -1, 64) +
		", " + strconv.FormatFloat(t.BestAsk, 'f', -1, 64) +
		", " + strconv.FormatFloat(t.BestAsk, 'f', -1, 64) + ")"
	_, err2 := db.Exec(query_insert)
	if err2 != nil {
		panic(err2.Error())
	}
}
func UpdateCandle(date string, h float64, l float64, c float64) {
	query_update := `update btc_jpy_live set ` +
		" high=" + strconv.FormatFloat(h, 'f', -1, 64) +
		", low=" + strconv.FormatFloat(l, 'f', -1, 64) +
		", close=" + strconv.FormatFloat(c, 'f', -1, 64) +
		" where date=" + `"` + date + `"`
	_, err3 := db.Exec(query_update)
	if err3 != nil {
		panic(err3.Error())
	}
}
func GetLatestCandle(date string) trade_def.BtcJpy {
	r, err := db.Query(`select * from btc_jpy_live where date="` + date + `" limit 1;`)
	if err != nil {
		panic(err.Error())
	}
	defer r.Close()
	var bj trade_def.BtcJpy
	for r.Next() {
		err = r.Scan(
			&bj.Date,
			&bj.Symbol,
			&bj.Open,
			&bj.High,
			&bj.Low,
			&bj.Close,
		)
		if err != nil {
			panic(err.Error())
		}
		break
	}
	return bj
}
func GetProfitList() ([]float64, []string) {
	rows, err := db.Query(`select * from btc_jpy_live_position where Profit is not NULL order by date;`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()

	var pos trade_def.Position
	var profits []float64
	var position_start_date []string
	for rows.Next() {
		err = rows.Scan(
			&pos.Date,
			&pos.Buy_or_sell,
			&pos.Price,
			&pos.Fix_date,
			&pos.Fix_price,
			&pos.Profit,
			&pos.Symbol,
		)
		if err != nil {
			panic(err.Error())
		}
		profits = append(profits, pos.Profit)
		position_start_date = append(position_start_date, pos.Date)
	}
	return profits, position_start_date

}
func GetCandleData() ([]trade_def.BtcJpy, float64, float64, int) {
	rows, err := db.Query(`select * from btc_jpy_live order by date;`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	var max float64
	var min float64
	max = 0.0
	min = 9999999.0
	records := make([]trade_def.BtcJpy, 0)
	for rows.Next() {
		var record trade_def.BtcJpy
		err = rows.Scan(
			&record.Date,
			&record.Symbol,
			&record.Open,
			&record.High,
			&record.Low,
			&record.Close,
		)
		if err != nil {
			panic(err.Error())
		}
		if max < record.High {
			max = record.High
		}
		if min > record.Low {
			min = record.Low
		}
		records = append(records, record)
	}

	return records, min, max, len(records)
}

func GetPositionData() []trade_def.Position {
	rows, err := db.Query(`select * from btc_jpy_live_position where Profit is not NULL order by date;`)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	records := make([]trade_def.Position, 0)
	for rows.Next() {
		var record trade_def.Position
		err = rows.Scan(
			&record.Date,
			&record.Buy_or_sell,
			&record.Price,
			&record.Fix_date,
			&record.Fix_price,
			&record.Profit,
			&record.Symbol,
		)
		if err != nil {
			panic(err.Error())
		}
		records = append(records, record)
	}

	return records
}
func GetProfitBacktest() float64 {
	var profit float64
	err := db.QueryRow(`select sum(profit) from btc_jpy_live_position_backtest;`).Scan(&profit)
	if err != nil {
		panic(err.Error())
	}
	return profit
}
