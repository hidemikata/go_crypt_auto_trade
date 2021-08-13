package trade_jadge_algo

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"fmt"
	"time"
)

type LongSma struct {
	sma_0 float64 //latest
	sma_1 float64
	sma_2 float64
	sma_3 float64
}

type ShortSma struct {
	sma_0 float64 //latest
	sma_1 float64
	sma_2 float64
}

type Sma struct {
	num_of_long  int
	num_of_short int
	Long         LongSma
	Short        ShortSma
	latest_min   float64
	latest_max   float64
	min_max_rate float64
}

var short_sma int
var long_sma int
var min_max_rate_number float64
var long_sma_margine int //解析は+3本必要
var now_date_margine int //現在時刻を省く
var latest_min_max_num int

func init() {
	short_sma = 5
	long_sma = 25
	min_max_rate_number = 0.005
	long_sma_margine = 3 //解析は+3本必要
	now_date_margine = 1 //現在時刻を省く
	latest_min_max_num = 5
}

func second_to_zero(t time.Time) string {
	min := fmt.Sprintf("%02d", t.Minute())
	h := fmt.Sprintf("%02d", t.Hour())
	d := fmt.Sprintf("%02d", t.Day())
	m := fmt.Sprintf("%02d", int(t.Month()))
	y := fmt.Sprintf("%02d", t.Year())
	return y + "-" + m + "-" + d + " " + h + ":" + min + ":00"
}

func NewSmaAlgorithm() *Sma {
	return &Sma{
		num_of_long:  long_sma,
		num_of_short: short_sma,
		min_max_rate: min_max_rate_number,
	}
}
func (sma_obj *Sma) Analisis(now time.Time) {
	margine_duration := time.Duration(now_date_margine)
	now_before_1min := now.Add(-(time.Minute * margine_duration))
	long_sma_duration := time.Duration(sma_obj.num_of_long + long_sma_margine + now_date_margine)
	now_before_long_sma := now.Add(-(time.Minute * long_sma_duration))
	latest_str := second_to_zero(now_before_1min)
	past_str := second_to_zero(now_before_long_sma)

	records := model.GetCandleBetweenDate(past_str, latest_str)

	l := LongSma{
		sma_0: calc_sma(records[:], sma_obj.num_of_long),
		sma_1: calc_sma(records[:len(records)-1], sma_obj.num_of_long),
		sma_2: calc_sma(records[:len(records)-2], sma_obj.num_of_long),
		sma_3: calc_sma(records[:len(records)-3], sma_obj.num_of_long),
	}
	s := ShortSma{
		sma_0: calc_sma(records[:], short_sma),
		sma_1: calc_sma(records[:len(records)-1], short_sma),
		sma_2: calc_sma(records[:len(records)-2], short_sma),
	}
	sma_obj.Long = l
	sma_obj.Short = s

	min, max := get_min_max(records[len(records)-latest_min_max_num : len(records)-1])

	sma_obj.latest_min = min
	sma_obj.latest_max = max
}
func get_min_max(records []trade_def.BtcJpy) (min float64, max float64) {
	min = 9999999.9
	max = 0.0
	for i := range records {
		if records[i].High > max {
			max = records[i].High
		}
		if records[i].Low < min {
			min = records[i].Low
		}
	}
	return min, max
}

func calc_sma(records []trade_def.BtcJpy, duration int) float64 {
	if len(records) < duration {
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

func (sma_obj *Sma) IsDbCollectedData(now time.Time) bool {
	num_of_collect := sma_obj.num_of_long + long_sma_margine + now_date_margine
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
		fmt.Println(now)
		panic("")
	}

	before_date := now.Add(-(time.Minute * num_of_duration))
	now_str := second_to_zero(now)
	before_date_str := second_to_zero(before_date)

	count := model.GetNumberOfCandleBetweenDate(before_date_str, now_str)

	return count-1 == num_of_collect //00秒〜00秒なので１個余分なので引く
}

func (sma_obj *Sma) IsTradeOrder() bool {
	if !check_sma(sma_obj) {
		//fmt.Println("sma ng")
		return false
	}
	if !check_rate_of_up(sma_obj) {
		//fmt.Println("rate of up ng")
		return false
	}

	return true
}
func check_rate_of_up(sma_obj *Sma) bool {
	rate := sma_obj.latest_min * sma_obj.min_max_rate
	//	fmt.Print(sma_obj.latest_max, sma_obj.latest_min, rate)
	//	fmt.Print((sma_obj.latest_max - sma_obj.latest_min))
	return (sma_obj.latest_max - sma_obj.latest_min) < rate
}
func check_sma(sma_obj *Sma) bool {
	if sma_obj.Short.sma_2 < sma_obj.Long.sma_2 &&
		sma_obj.Short.sma_1 < sma_obj.Long.sma_1 &&
		sma_obj.Short.sma_0 > sma_obj.Long.sma_0 &&
		sma_obj.Short.sma_2 < sma_obj.Short.sma_1 &&
		sma_obj.Short.sma_1 < sma_obj.Short.sma_0 &&
		sma_obj.Long.sma_3 < sma_obj.Long.sma_2 &&
		sma_obj.Long.sma_2 < sma_obj.Long.sma_1 &&
		sma_obj.Long.sma_1 < sma_obj.Long.sma_0 {
		return true
	}
	return false
}

func (sma_obj *Sma) IsTradeFix() bool {
	return sma_obj.Short.sma_0 < sma_obj.Long.sma_0
}

func (sma_obj *Sma) SetParam(sma ...int) {
	fmt.Println("sma set param ", sma)
	sma_obj.num_of_long = sma[0]
	sma_obj.num_of_short = sma[1]
	sma_obj.min_max_rate = float64(sma[2]) / 1000
}
