package long_sma_algo

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/app/trade_jadge_algo"
	"btcanallive_refact/config"
	"fmt"
	"time"
)

type LongLongSma struct {
	sma_0 float64 //latest
}

type LongSma struct {
	num_of_long_long int
	LongLong         LongLongSma
	latest_candle    trade_def.BtcJpy //直近のローソク足
	past_candle      trade_def.BtcJpy //最も古いのローソク足
}

var sma_margine int      //解析は+3本必要
var now_date_margine int //現在時刻を省く

func init() {
	sma_margine = 3      //解析は+3本必要
	now_date_margine = 1 //現在時刻を省く
}

func second_to_zero(t time.Time) string {
	min := fmt.Sprintf("%02d", t.Minute())
	h := fmt.Sprintf("%02d", t.Hour())
	d := fmt.Sprintf("%02d", t.Day())
	m := fmt.Sprintf("%02d", int(t.Month()))
	y := fmt.Sprintf("%02d", t.Year())
	return y + "-" + m + "-" + d + " " + h + ":" + min + ":00"
}

func NewLongSmaAlgorithm() *LongSma {
	return &LongSma{
		num_of_long_long: config.Config.SmaLongLong,
	}
}
func (sma_obj *LongSma) Analisis(now time.Time) {
	margine_duration := time.Duration(now_date_margine)
	now_before_1min := now.Add(-(time.Minute * margine_duration))
	long_sma_duration := time.Duration(sma_obj.num_of_long_long + sma_margine + now_date_margine)
	now_before_long_sma := now.Add(-(time.Minute * long_sma_duration))
	latest_str := second_to_zero(now_before_1min)
	past_str := second_to_zero(now_before_long_sma)

	records := model.GetCandleBetweenDate(past_str, latest_str)

	ll := LongLongSma{
		sma_0: calc_sma(records[:], sma_obj.num_of_long_long),
	}
	sma_obj.LongLong = ll

	sma_obj.latest_candle = records[len(records)-1]
	sma_obj.past_candle = records[0]
}

func calc_sma(records []trade_def.BtcJpy, duration int) float64 {
	if len(records) < duration {
		fmt.Println("***long long records =", len(records))
	}
	total := 0.0
	for i := range records {
		total += records[i].Close
	}
	return total / float64(len(records))
}

func (sma_obj *LongSma) IsDbCollectedData(now time.Time) bool {
	return true
}

func (sma_obj *LongSma) IsTradeOrder() bool {
	if !check_sma(sma_obj) {
		return false
	}

	if !check_price_past_and_latest(sma_obj) {
		return false
	}
	if !check_in_range(sma_obj) {
		return false
	}

	return true
}
func check_in_range(sma_obj *LongSma) bool {
	return sma_obj.latest_candle.Low-sma_obj.LongLong.sma_0 < sma_obj.latest_candle.Low*0.005
}

func check_sma(sma_obj *LongSma) bool {
	return sma_obj.latest_candle.Low > sma_obj.LongLong.sma_0
}
func check_price_past_and_latest(sma_obj *LongSma) bool {
	return sma_obj.latest_candle.Close > sma_obj.past_candle.Close
}

func (sma_obj *LongSma) IsTradeFix() bool {
	return false
}

func (sma_obj *LongSma) BacktestSetParam(params []int) {
}
func (sma_obj *LongSma) FixRealTick(t trade_def.Ticker) bool {
	return false
}

func (sma_obj *LongSma) CreateBacktestParams() []trade_jadge_algo.BacktestParams {

	p := make([]trade_jadge_algo.BacktestParams, 0)
	return p
}

func (sma_obj *LongSma) Name() string {
	return "long sma up rate."
}
