package http_template

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
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

	candle_data := model.GetCandleData()

	if err := t.Execute(w, struct {
		Title     string
		Message   string
		Time      time.Time
		Profit    []Data
		CanleDate []trade_def.BtcJpy
	}{
		Title:     title,
		Message:   "こんにちは！",
		Time:      time.Now(),
		Profit:    d,
		CanleDate: candle_data,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
