package trade_jadge_algo

import (
	"btcanallive_refact/app/trade_def"
	"time"
)

type BacktestParams struct {
	Ti     TradeInterface
	Params []int
}

type TradeInterface interface {
	IsDbCollectedData(time.Time) bool       //データが全部揃ってるか
	Analisis(time.Time)                     //解析
	IsTradeOrder() bool                     //仕掛けるか
	IsTradeFix() bool                       //決済するかどうか
	BacktestSetParam([]int)                 //パラメータを設定する。
	FixRealTick(t trade_def.Ticker) bool    //tickレベルで決済
	CreateBacktestParams() []BacktestParams //バックテストのパラメータ生成
}
