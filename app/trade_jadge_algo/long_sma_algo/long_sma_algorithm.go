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
	sma_1 float64
	sma_2 float64
	sma_3 float64
}

type LongSma struct {
	num_of_long_long int
	LongLong         LongLongSma
	latest_candle    trade_def.BtcJpy //直近のローソク足
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
		sma_1: calc_sma(records[:len(records)-1], sma_obj.num_of_long_long),
		sma_2: calc_sma(records[:len(records)-2], sma_obj.num_of_long_long),
		sma_3: calc_sma(records[:len(records)-3], sma_obj.num_of_long_long),
	}
	sma_obj.LongLong = ll

	sma_obj.latest_candle = records[len(records)-1]
}

func calc_sma(records []trade_def.BtcJpy, duration int) float64 {
	if len(records) < duration {
		fmt.Println(len(records), duration)
		fmt.Println(time.Now())
		panic("")
	}
	total := 0.0
	start_i := len(records) - duration
	record_latest := records[start_i:]
	for i := range record_latest {
		total += record_latest[i].Close
	}
	return total / float64(len(record_latest))
}

func (sma_obj *LongSma) IsDbCollectedData(now time.Time) bool {
	num_of_collect := sma_obj.num_of_long_long + sma_margine + now_date_margine
	num_of_duration := time.Duration(num_of_collect)

	if now.Second() > 50 {
		//50秒から５９秒の間は待つ
		sleep_times := 0
		for {
			//tikcerのタイムで足確定判断して、PCのタイムで解析スタートしようとしているので時間がずれる
			fmt.Println(now.Second(), "is_db_collected_data time is not 00 sec sleep0.5...")
			time.Sleep(500 * time.Millisecond) //now時間がX分59秒になることがあるので0.5秒待つ
			sleep_times++
			now = time.Now()
			if now.Second() == 0 {
				break
			}
		}
	} else if now.Second() > 10 {
		//10秒以上差が出てたら落とす
		fmt.Println(now, now.Second())
		panic("")
	}

	before_date := now.Add(-(time.Minute * num_of_duration))
	now_str := second_to_zero(now)
	before_date_str := second_to_zero(before_date)

	count := model.GetNumberOfCandleBetweenDate(before_date_str, now_str)
	return count-1 == num_of_collect //00秒〜00秒なので１個余分なので引く
}

func (sma_obj *LongSma) IsTradeOrder() bool {
	if !check_sma(sma_obj) {
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
	if sma_obj.LongLong.sma_2 < sma_obj.LongLong.sma_1 &&
		sma_obj.LongLong.sma_1 < sma_obj.LongLong.sma_0 &&
		sma_obj.latest_candle.Low < sma_obj.LongLong.sma_0 {
		return true
	}
	return false
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