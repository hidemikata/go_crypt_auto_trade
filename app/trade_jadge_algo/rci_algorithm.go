package trade_jadge_algo

import "time"

type Rci struct {
}

func NewRciAlgorithm() *Rci {
	return &Rci{}
}

func (obj *Rci) IsDbCollectedData(now time.Time) bool {
	return true
}
func (obj *Rci) Analisis(anal_time time.Time) {

}
func (obj *Rci) IsTradeOrder() bool {
	return true
}
func (obj *Rci) IsTradeFix() bool {
	return false

}

func (sma_obj *Rci) GetLatestAskBid() (float64, float64) {
	return 0.0, 0.0
}
