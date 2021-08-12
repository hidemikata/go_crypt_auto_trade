package trade_def

type Position struct {
	Date        string
	Buy_or_sell string
	Price       float64
	Fix_date    string
	Fix_price   float64
	Profit      float64
	Symbol      string
}
type BtcJpy struct {
	Date   string
	Symbol string
	Open   float64
	High   float64
	Low    float64
	Close  float64
}

type Ticker struct {
	ProductCode     string  `json:"product_code"`
	State           string  `json:"state"`
	Timestamp       string  `json:"timestamp"`
	TickID          int     `json:"tick_id"`
	BestBid         float64 `json:"best_bid"`
	BestAsk         float64 `json:"best_ask"` //買い
	BestBidSize     float64 `json:"best_bid_size"`
	BestAskSize     float64 `json:"best_ask_size"`
	TotalBidDepth   float64 `json:"total_bid_depth"`
	TotalAskDepth   float64 `json:"total_ask_depth"`
	MarketBidSize   int     `json:"market_bid_size"`
	MarketAskSize   int     `json:"market_ask_size"`
	Ltp             int     `json:"ltp"`
	Volume          float64 `json:"volume"`
	VolumeByProduct float64 `json:"volume_by_product"`
}
