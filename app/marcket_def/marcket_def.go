package marcket_def
import(
	"btcanallive_refact/app/trade_def"//$B$d$j$?$/$J$$(B
)
type Marcket interface {
	StartGettingRealTimeTicker(chan<- trade_def.Ticker)
	PutOrder()
	FixOrder()
	GetTicker() trade_def.Ticker
}


