package trade_jadge_algo

type TradeInterface interface {
	IsDbCollectedData() bool
	Analisis()
	IsTradeOrder() bool
	IsTradeFix() bool
}

