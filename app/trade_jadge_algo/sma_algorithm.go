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
}

var short_sma int
var long_sma int
var long_sma_margine int //$B2r@O$O(B+$B#2K\I,MW(B
var now_date_margine int //$B8=:_;~9o$r>J$/(B

func init() {
	short_sma = 5
	long_sma = 25
	long_sma_margine = 3 //$B2r@O$O(B+$B#2K\I,MW(B
	now_date_margine = 1 //$B8=:_;~9o$r>J$/(B
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
	}
}
func (sma_obj *Sma) Analisis() {
	now := time.Now()
	margine_duration := time.Duration(now_date_margine)
	now_before_1min := now.Add(-(time.Minute * margine_duration))
	long_sma_duration := time.Duration(long_sma + long_sma_margine + now_date_margine)
	now_before_long_sma := now.Add(-(time.Minute * long_sma_duration))
	latest_str := second_to_zero(now_before_1min)
	past_str := second_to_zero(now_before_long_sma)

	records := model.GetCandleBetweenDate(past_str, latest_str)

	l := LongSma{
		sma_0: calc_sma(records[:], long_sma),
		sma_1: calc_sma(records[:len(records)-1], long_sma),
		sma_2: calc_sma(records[:len(records)-2], long_sma),
		sma_3: calc_sma(records[:len(records)-3], long_sma),
	}
	s := ShortSma{
		sma_0: calc_sma(records[:], short_sma),
		sma_1: calc_sma(records[:len(records)-1], short_sma),
		sma_2: calc_sma(records[:len(records)-2], short_sma),
	}
	sma_obj.Long = l
	sma_obj.Short = s
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

func (sma_obj *Sma) IsDbCollectedData() bool {
	num_of_collect := long_sma + long_sma_margine + now_date_margine
	num_of_duration := time.Duration(num_of_collect)
	now := time.Now()

	if now.Second() > 50 {
		//50$BIC$+$i#5#9IC$N4V$OBT$D(B
		sleep_times := 0
		for {
			//tikcer$B$N%?%$%`$GB-3NDjH=CG$7$F!"(BPC$B$N%?%$%`$G2r@O%9%?!<%H$7$h$&$H$7$F$$$k$N$G;~4V$,$:$l$k(B
			fmt.Println(now.Second(), "is_db_collected_data time is not 00 sec sleep0.5...")
			time.Sleep(500 * time.Millisecond) //now$B;~4V$,(BX$BJ,(B59$BIC$K$J$k$3$H$,$"$k$N$G(B0.5$BICBT$D(B
			sleep_times++
			now = time.Now()
			if now.Second() == 0 {
				break
			}
		}
	} else if now.Second() > 10 {
		//10$BIC0J>e:9$,=P$F$?$iMn$H$9(B
		fmt.Println(now)
		panic("")
	}

	before_date := now.Add(-(time.Minute * num_of_duration))
	now_str := second_to_zero(now)
	before_date_str := second_to_zero(before_date)

	count := model.GetNumberOfCandleBetweenDate(before_date_str, now_str)

	return count-1 == num_of_collect //00$BIC!A(B00$BIC$J$N$G#18DM>J,$J$N$G0z$/(B
}

func (sma_obj *Sma) IsTradeOrder() bool {
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
