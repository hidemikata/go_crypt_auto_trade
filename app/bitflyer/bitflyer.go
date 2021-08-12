package bitflyer

import (
	"btcanallive_refact/app/marcket_def"
	"btcanallive_refact/app/trade_def"
	"btcanallive_refact/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const BasePath = "https://api.bitflyer.com"
const TickerUrl = "/v1/getticker"
const MyAllOrderCancel = "/v1/me/cancelallchildorders"
const MySendOrder = "/v1/me/sendchildorder"

type JsonRPC2 struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Result  interface{} `json:"result,omitempty"`
	Id      *int        `json:"id,omitempty"`
}

type SubscribeParams struct {
	Channel string `json:"channel"`
}

type Bitflyer struct {
	ProductCode string
}

func NewBitflyer(product_code string) marcket_def.Marcket {
	return &Bitflyer{
		ProductCode: product_code,
	}
}

func setPrivateHeader(req *http.Request) {
	req.Header.Add("ACCESS-KEY", config.Config.ApiKey)
	t := time.Now().Unix()
	ts := strconv.FormatInt(t, 10)
	req.Header.Add("ACCESS-TIMESTAMP", ts)
	req_body, _ := req.GetBody()
	body_byte, _ := io.ReadAll(req_body)
	var hira = ts + req.Method + req.URL.RequestURI() + string(body_byte)
	mac := hmac.New(sha256.New, []byte(config.Config.ApiSecret))
	mac.Write([]byte(hira))
	req.Header.Add("ACCESS-SIGN", hex.EncodeToString(mac.Sum(nil)))
	req.Header.Add("Content-Type", "application/json")
}

func AllOrderCancel() {
	fmt.Println("AllOrderCancel")
	client := &http.Client{CheckRedirect: nil}

	req_body := strings.NewReader(`{"product_code":"` + "BTC_JPY" + `"}`)
	req, _ := http.NewRequest("POST", BasePath+MyAllOrderCancel, req_body)
	setPrivateHeader(req)

	res, _ := client.Do(req)
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
}

func (api *Bitflyer) GetTicker() trade_def.Ticker {
	fmt.Println("GetTicker")
	resp, _ := http.Get(BasePath + TickerUrl + "?product_code=BTC_JPY")
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var decoded trade_def.Ticker
	json.Unmarshal(body, &decoded)
	return decoded
}

func (api *Bitflyer) PutOrder() {
	fmt.Println("put order to bitflyer api")
}

func (api *Bitflyer) FixOrder() {
	fmt.Println("fix order to bitflyer api")
}

func (api *Bitflyer) StartGettingRealTimeTicker(ch chan<- trade_def.Ticker) {
	u := url.URL{Scheme: "wss", Host: "ws.lightstream.bitflyer.com", Path: "/json-rpc"}
	fmt.Println("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("dial:", err)
	}
	defer c.Close()

	channel := fmt.Sprintf("lightning_ticker_%s", api.ProductCode)
	if err := c.WriteJSON(&JsonRPC2{Version: "2.0", Method: "subscribe", Params: &SubscribeParams{channel}}); err != nil {
		fmt.Println("subscribe:", err)
		return
	}
OUTER:
	for {
		message := new(JsonRPC2)
		if err := c.ReadJSON(message); err != nil {
			fmt.Println("read:", err)
			return
		}

		if message.Method == "channelMessage" {
			switch v := message.Params.(type) {
			case map[string]interface{}:
				for key, binary := range v {
					if key == "message" {
						marshaTic, err := json.Marshal(binary)
						if err != nil {
							fmt.Println("err1")
							continue OUTER
						}
						var ticker trade_def.Ticker
						if err := json.Unmarshal(marshaTic, &ticker); err != nil {
							fmt.Println("err2")
							continue OUTER
						}
						ch <- ticker
					}
				}
			}
		}
	}
}
