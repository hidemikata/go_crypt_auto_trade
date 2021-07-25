package marcket_def
import(
	"btcanallive_refact/app/trade_def"//やりたくない
)
type Marcket interface {
	StartGettingRealTimeTicker(chan<- trade_def.Ticker)
	PutOrder()
	FixOrder()
	GetTicker() trade_def.Ticker
}


