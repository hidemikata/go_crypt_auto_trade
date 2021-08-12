package trade_jadge_algo

type Rci struct {
}

func NewRciAlgorithm() *Rci {
	return &Rci{}
}

func (obj *Rci) IsDbCollectedData() bool {
	return true
}
func (obj *Rci) Analisis() {

}
func (obj *Rci) IsTradeOrder() bool {
	return true
}
func (obj *Rci) IsTradeFix() bool {
	return true

}
