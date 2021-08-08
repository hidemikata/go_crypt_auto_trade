package http_template

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Number     int
	ProfitSum  float64
	Profit     float64
	DateSecond int
}

func checkAuth(r *http.Request) bool {
	id, pass, ok := r.BasicAuth()
	if !ok {
		return false
	}
	return id == "bakueki" && pass == "aba"
}

var layout = "2006-01-02 15:04:05"

func stringToTime(str string) time.Time {
	t, _ := time.Parse(layout, str)
	return t
}
func diff_second(t time.Time) int {
	day1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	day2 := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
	duration := day2.Sub(day1)
	return int(duration.Seconds())
}
func get_second_from_00(str string) int {
	t := stringToTime(str)
	return diff_second(t)
}

func ProfitView(w http.ResponseWriter, r *http.Request) {

	if !checkAuth(r) {
		w.Header().Add("WWW-Authenticate", `Basic realm="my private area"`)
		w.WriteHeader(http.StatusUnauthorized) // 401コード
		// 認証失敗時の出力内容
		w.Write([]byte("401 認証失敗\n"))
		return
	}

	t, err := template.ParseFiles("app/http_template/profit_view.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}

	profits, position_date := model.GetProfitList()
	var d []Data
	var profit_sum float64
	for i, v := range profits {
		profit_sum += v
		d = append(d, Data{i, profit_sum, v, get_second_from_00(position_date[i])})
	}

	var title string
	if profit_sum > 0 {
		title = "爆益"
	} else {
		title = "爆損"
	}

	candle_data, candle_min, candle_max := model.GetCandleData()

	sma1 := make([]float64, 0)
	sma2 := make([]float64, 0)

	for i := range candle_data {
		sma1 = append(sma1, calc_sma(candle_data[:i], 5))
		sma2 = append(sma2, calc_sma(candle_data[:i], 25))
	}
	var candle_year []string
	var candle_month []string
	var candle_day []string
	for _, v := range candle_data {
		t, _ := time.Parse(layout, v.Date)

		candle_year = append(candle_year, fmt.Sprintf("%d", t.Year()))
		candle_month = append(candle_month, fmt.Sprintf("%d", int(t.Month())))
		candle_day = append(candle_day, fmt.Sprintf("%d", t.Day()))

	}
	if err := t.Execute(w, struct {
		Title       string
		Message     string
		Time        time.Time
		Profit      []Data
		CanleDate   []trade_def.BtcJpy
		CandleMax   float64
		CandleMin   float64
		Sma1        []float64
		Sma2        []float64
		CandleYear  []string
		CandleMonth []string
		CandleDay   []string
	}{
		Title:       title,
		Message:     "こんにちは！",
		Time:        time.Now(),
		Profit:      d,
		CanleDate:   candle_data,
		CandleMax:   candle_max,
		CandleMin:   candle_min,
		Sma1:        sma1,
		Sma2:        sma2,
		CandleYear:  candle_year,
		CandleMonth: candle_month,
		CandleDay:   candle_day,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func calc_sma(records []trade_def.BtcJpy, duration int) float64 {
	if len(records) < duration {
		return 0.0
	}
	total := 0.0
	start_i := len(records) - duration
	record_latest := records[start_i:]
	for i := range record_latest {
		total += record_latest[i].Close
	}
	return total / float64(len(record_latest))
}
