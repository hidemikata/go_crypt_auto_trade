package http_template

import (
	"btcanallive_refact/app/model"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Data struct {
	Number int
	Profit float64
}

func ProfitView(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("app/http_template/profit_view.html")
	if err != nil {
		log.Fatalf("template error: %v", err)
	}
	profits := model.GetProfitList()
	var d []Data

	for i, v := range profits {
		d = append(d, Data{i, v})
	}
	fmt.Println(d)

	if err := t.Execute(w, struct {
		Title   string
		Message string
		Time    time.Time
		Profit  []Data
	}{
		Title:   "爆損",
		Message: "こんにちは！",
		Time:    time.Now(),
		Profit:  d,
	}); err != nil {
		log.Printf("failed to execute template: %v", err)
	}
}
