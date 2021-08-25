package trade_jadge_algo

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/config"
	"fmt"
	"math"
	"time"
)

type Rci struct {
	rci_long     int
	rci_buy_rate float64 //これ以下なら買う-100から100
	rci          float64
}

var rci_now_date_margine int //現在時刻を省く

func init() {
	rci_now_date_margine = 1 //現在時刻を省く
}

func NewRciAlgorithm() *Rci {
	return &Rci{
		rci_long:     config.Config.RciLong,
		rci_buy_rate: config.Config.RciRate,
		rci:          0.0,
	}
}

func (obj *Rci) Analisis(anal_time time.Time) {
	margine_duration := time.Duration(rci_now_date_margine)
	now_before_1min := anal_time.Add(-(time.Minute * margine_duration))
	long_duration := time.Duration(obj.rci_long)
	now_before_long := anal_time.Add(-(time.Minute * long_duration))
	latest_str := second_to_zero(now_before_1min)
	past_str := second_to_zero(now_before_long)

	records := model.GetCandleBetweenDate(past_str, latest_str)
	tmp := make([]trade_def.BtcJpy, len(records))

	copy(tmp, records)

	price_order := make([]trade_def.BtcJpy, 0)

	for j := 0; j < len(records); j++ {
		var top_price_data trade_def.BtcJpy
		var tmp_i int
		for i, v := range tmp {
			if v.Close >= top_price_data.Close {
				top_price_data = v
				tmp_i = i
			}
		}
		price_order = append(price_order, top_price_data)
		tmp = remove(tmp, tmp_i)
	}

	var date_price_square float64
	var hizuke_kakaku_sa float64
	for i_r, v_r := range records {
		for i_po, v_po := range price_order {
			if v_po.Date == v_r.Date {
				hizuke_kakaku_sa = float64((len(records) - i_r) - (i_po + 1))
				date_price_square = date_price_square + math.Abs(hizuke_kakaku_sa*hizuke_kakaku_sa)
				break
			}
		}
	}
	bunbo := (len(records)*len(records) - 1) * len(records)
	rci_tmp := float64(date_price_square) * 6 / float64(bunbo)
	rci := (1 - rci_tmp) * 100
	obj.rci = rci

}

func remove(slice []trade_def.BtcJpy, s int) []trade_def.BtcJpy {
	return append(slice[:s], slice[s+1:]...)
}

func (obj *Rci) IsTradeOrder() bool {
	if obj.rci_long == 0 {
		return true
	}
	if obj.rci < obj.rci_buy_rate {
		return true
	}
	//fmt.Println(" rci ng:", obj.rci)
	return false
}

func (obj *Rci) IsTradeFix() bool {
	return false

}

func (obj *Rci) IsDbCollectedData(now time.Time) bool {
	num_of_collect := obj.rci_long + rci_now_date_margine
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

func (obj *Rci) SetParam(rci ...int) {
	fmt.Println("rci set param ", rci)
	obj.rci_long = rci[0]
	obj.rci_buy_rate = float64(rci[1])
}
