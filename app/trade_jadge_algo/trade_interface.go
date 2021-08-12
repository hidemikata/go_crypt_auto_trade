package trade_jadge_algo

import "time"

type TradeInterface interface {
	IsDbCollectedData(time.Time) bool    //データが全部揃ってるか
	Analisis(time.Time)                  //解析
	IsTradeOrder() bool                  //仕掛けるか
	IsTradeFix() bool                    //けっさいするか
	GetLatestAskBid() (float64, float64) //解析した時の売りねとかいねを取得。本当はここじゃない
}
