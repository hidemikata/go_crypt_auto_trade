package web_page

import (
	"btcanallive_refact/app/http_template"
	"net/http"
)

func StartWebServer() {
	http.HandleFunc("/profit", http_template.ProfitView)
	http.ListenAndServe(":8080", nil)
}
