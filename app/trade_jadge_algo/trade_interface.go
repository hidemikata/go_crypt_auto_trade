package trade_jadge_algo

type TradeInterface interface {
	IsDbCollectedData() bool //$B%G!<%?$,A4ItB7$C$F$k$+(B
	Analisis()               //$B2r@O(B
	IsTradeOrder() bool      //$B;E3]$1$k$+(B
	IsTradeFix() bool        //$B$1$C$5$$$9$k$+(B
}
