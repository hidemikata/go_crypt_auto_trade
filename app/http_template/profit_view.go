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
func checkAuth(r *http.Request) bool {
    id, pass, ok := r.BasicAuth()
    if ok == false{
        return false
    }
    return id == "bakueki" && pass == "aba"
}
func ProfitView(w http.ResponseWriter, r *http.Request) {
    if checkAuth(r) == false{
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
