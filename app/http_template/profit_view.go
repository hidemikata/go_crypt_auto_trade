package http_template

import (
	"btcanallive_refact/app/model"
	"btcanallive_refact/app/trade_def"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

//label
type PositionView struct {
	Id         int
	ProfitText string
	LabelColor string
	LabelSize  int
	Position   string
	Shape      string
}

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

	candle_data, _, _, _ := model.GetCandleData()
	positions := model.GetPositionData()

	position_time := make([]PositionView, 0)
	for _, v_pos := range positions {
		for i, v_candle := range candle_data {
			color := "blue"
			if second_to_zero(timeComvert(v_pos.Date)) == second_to_zero(timeComvert(v_candle.Date)) {
				if v_pos.Profit < 0 {
					color = "red"
				}
				position_tmp := PositionView{
					Id:         i,
					ProfitText: strconv.FormatFloat(v_pos.Profit, 'f', -1, 64),
					LabelColor: color,
					LabelSize:  10,
					Position:   "aboveBar",
					Shape:      "arrowDown",
				}
				position_time = append(position_time, position_tmp)
			} else if second_to_zero(timeComvert(v_pos.Fix_date)) == second_to_zero(timeComvert(v_candle.Date)) {
				position_tmp := PositionView{
					Id:         i,
					ProfitText: "",
					LabelColor: "black",
					LabelSize:  5,
					Position:   "belowBar",
					Shape:      "arrowUp",
				}
				position_time = append(position_time, position_tmp)
				break

			}
		}
	}

	if err := t.Execute(w, struct {
		Title        string
		Message      string
		Time         time.Time
		Profit       []Data
		CanleData    []trade_def.BtcJpy
		PositionTime []PositionView
	}{
		Title:        title,
		Message:      "こんにちは！",
		Time:         time.Now(),
		Profit:       d,
		CanleData:    candle_data,
		PositionTime: position_time,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}

func timeComvert(date string) time.Time {
	t, _ := time.Parse(layout, date)
	return t
}
func second_to_zero(t time.Time) string {
	min := fmt.Sprintf("%02d", t.Minute())
	h := fmt.Sprintf("%02d", t.Hour())
	d := fmt.Sprintf("%02d", t.Day())
	m := fmt.Sprintf("%02d", int(t.Month()))
	y := fmt.Sprintf("%02d", t.Year())
	return y + "-" + m + "-" + d + " " + h + ":" + min + ":00"
}
