package trade_jadge_algo

type TradeInterface interface {
	IsDbCollectedData() bool //データが全部揃ってるか
	Analisis()               //解析
	IsTradeOrder() bool      //仕掛けるか
	IsTradeFix() bool        //けっさいするか
}
