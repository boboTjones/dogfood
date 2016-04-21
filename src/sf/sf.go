package sf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"util"
)

type Order struct {
	OK        bool    `json:"ok"`
	Account   string  `json:"account"`
	Venue     string  `json:"venue"`
	Symbol    string  `json:"stock"`
	Price     float64 `json:"price"`
	Qty       float64 `json:"qty"`
	Direction string  `json:"direction"`
	Type      string  `json:"orderType"`
}

type Restart struct {
	InstanceId float64  `json:"instanceId"`
	Account    string   `json:"account"`
	Tickers    []string `json:"tickers"`
	Venues     []string `json:"venues"`
}

type Quote struct {
	OK       bool    `json:"ok"`
	Symbol   string  `json:"symbol"`
	Venue    string  `json:"venue"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
	BidSize  float64 `json:"bidSize"`
	AskSize  float64 `json:"askSize"`
	BidDepth float64 `json:"bidDepth"`
	AskDept  float64 `json:"askDepth"`
	Last     float64 `json:"last"`
	LastSize float64 `json:"lastSize"`
}

type LevelInfo struct {
	OK         bool     `json:"ok"`
	InstanceId float64  `json:"instanceId"`
	Account    string   `json:"account"`
	Tickers    []string `json:"tickers"`
	Venues     []string `json:"venues"`
	Seconds    float64  `json:"secondsPerTradingDay"`
}

type Stocks struct {
	OK      bool                `json:"ok"`
	Symbols []map[string]string `json:"symbols"`
}

type Bid struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
	IsBuy bool    `json:"isBuy"`
}

type Ask struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
	IsBuy bool    `json:"isBuy"`
}

type Stock struct {
	OK     bool   `json:"ok"`
	Venue  string `json:"venue"`
	Symbol string `json:"symbol"`
	Stamp  string `json:"ts"`
	Bids   []Bid  `json:"bids"`
	Asks   []Ask  `json:"asks"`
}

type Fill struct {
	Price float64 `json:"price"`
	Qty   float64 `json:"qty"`
	Stamp string  `json:"ts"`
}

type Trade struct {
	OK        bool    `json:"ok"`
	Venue     string  `json:"venue"`
	Symbol    string  `json:"symbol"`
	Direction string  `json:"direction"`
	Qty       float64 `json:"originalQty"`
	Unfilled  float64 `json:"qty"`
	Price     float64 `json:"price"`
	Type      string  `json:"orderType"`
	Id        float64 `json:"id"`
	Account   string  `json:"account"`
	Stamp     string  `json:"ts"`
	Fills     []Fill  `json:"fills"`
	Filled    float64 `json:"totalFilled"`
	Open      bool    `json:"open"`
}

type Trades struct {
	OK     bool    `json:"ok"`
	Orders []Trade `json:"orders"`
}

func Heartbeat(t *url.URL, key string) bool {
	//get /venues/:venue/heartbeat
	var res map[string]interface{}
	t.Path = "/ob/api/heartbeat"
	r, _ := util.Get(t, key)
	json.Unmarshal(r.Body, &res)
	return res["ok"].(bool)
}

func Vbeat(t *url.URL, v, key string) bool {
	//get /venues/:venue/heartbeat
	var res map[string]interface{}
	t.Path = "/ob/api/venues/" + v + "/heartbeat"
	r, _ := util.Get(t, key)
	json.Unmarshal(r.Body, &res)
	return res["ok"].(bool)
}

func GetStocks(t *url.URL, v, key string) (Stocks, error) {
	var stocks Stocks
	//get /venues/:venue/stocks
	t.Path = "/ob/api/venues/" + v + "/stocks"
	r, err := util.Get(t, key)
	if err = json.Unmarshal(r.Body, &stocks); err != nil {
		return stocks, err
	}
	return stocks, err
}

func StockInfo(t *url.URL, v, s, key string) (Stock, error) {
	//get /venues/:venue/stocks/:stock
	var infos Stock
	t.Path = "/ob/api/venues/" + v + "/stocks/" + s
	r, err := util.Get(t, key)
	if err = json.Unmarshal(r.Body, &infos); err != nil {
		return infos, err
	}
	return infos, err
}

func GetQuote(t *url.URL, v, s, key string) (Quote, error) {
	//get /venues/:venue/stocks/:stock/quote
	var quote Quote
	var err error

	t.Path = "/ob/api/venues/" + v + "/stocks/" + s + "/quote"
	r, err := util.Get(t, key)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println(string(r.Body))

	if err := json.Unmarshal(r.Body, &quote); err != nil {
		fmt.Println(err)
		return quote, err
	}
	return quote, err
}

func MakeOrder(t *url.URL, data []byte, v, s, key string) (*util.Response, error) {
	//post /venues/:venue/stocks/:stock/orders
	t.Path = "/ob/api/venues/" + v + "/stocks/" + s + "/orders"
	r, err := util.Post(t, data, key)
	return r, err
}

func ShowOrder(v, s, id string) (string, string) {
	//get /venues/:venue/stocks/:stock/orders/:id
	return "GET", "/ob/api/venues/" + v + "/stocks/" + s + "/orders/" + id

}

func Cancel(t *url.URL, v, s, key string, id float64) {
	//delete /venues/:venue/stocks/:stock/orders/:order
	t.Path = "/ob/api/venues/" + v + "/stocks/" + s + "/orders/" + fmt.Sprintf("%g", id)
	r, err := util.Del(t, key)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(r.Body))
}

func GetOrdersForAcct(v, a string) (string, string) {
	//get /venues/:venue/accounts/:account/orders
	return "GET", "/ob/api/venues/" + v + "/accounts/" + a + "/orders"

}

func Mine(t *url.URL, v, a, s, key string) (Trades, error) {
	var trades Trades
	//get /venues/:venue/accounts/:account/stocks/:stock/orders
	t.Path = "/ob/api/venues/" + v + "/accounts/" + a + "/stocks/" + s + "/orders"
	r, err := util.Get(t, key)
	fmt.Println(string(r.Body))
	if err = json.Unmarshal(r.Body, &trades); err != nil {
		return trades, err
	}
	return trades, err

}

func NewOrder(a, v, s, d, ot string, p, qty float64) *Order {
	return &Order{
		Account:   a,
		Venue:     v,
		Symbol:    s,
		Direction: d,
		Type:      ot,
		Price:     p,
		Qty:       qty,
	}
}

func cookieMonster(filename string) (map[string]string, float64) {
	cooks := make(map[string]string)
	var instanceId float64
	var levels map[string]interface{}

	x := util.SlurpFromFile(filename)
	b := bytes.Split(x, []byte("\r\n"))
	for _, v := range b {
		c := bytes.Split(v, []byte("\t"))
		l := len(c)
		if l > 1 {
			cooks[string(c[l-2])] = string(c[l-1])
		}
	}
	if err := json.Unmarshal([]byte(cooks["levelInstances"]), &levels); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// Get the last one. Yuck.
	for _, v := range levels {
		//level = k
		instanceId = v.(float64)
	}
	return cooks, instanceId
}

func Login(u, p string) ([]*http.Cookie, error) {
	stupid := `RFplNmpTazFRWFoxZFJ2SnlxNytmVEVlZkthL21ra2VhU3JZT0ZKVHRVTGE0Q1EzendCV3d3UnJidm1mcyszaGkrVDh6aWppZTJHWnZGcHZtU0lxSVZGVEc2REFNNFg5YytxVjh0Uk80M0hFb3hrTFlRd2Zzb2oxaVZOL1dxK0YzbExLaURicEtFaXVXTHJFUWV2QTNRPT0tLTdvSGREZVhXcTVLcU5LMXBlZHh4eFE9PQ%3D%3D--5d724eb099f745d8ae33003daf1191346fd25281`
	d := []byte(fmt.Sprintf("session[username]=%s&session[password]=%s", u, p))
	t, _ := url.Parse("https://www.stockfighter.io/ui/login")
	req, err := http.NewRequest("POST", t.String(), bytes.NewReader(d))
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{
		Name:  "_playerui_session",
		Value: stupid,
	})
	req.Header.Set("X-CSRF-Token", `4l+3kcK2d/p0reAC/+lxCbTPHGrtB6m4Z/FdgYF4BgtV/im1RzcbOPaNLCUPEjrx++aN5QA04jjyq/h+6ty7GQ==`)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Cookies(), err
}

func RestartLastLevel(cooks []*http.Cookie, instanceId float64) (Restart, error) {
	var infos Restart

	u, _ := url.Parse(fmt.Sprintf("https://www.stockfighter.io/gm/instances/%g/restart", instanceId))
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return infos, err
	}

	// invalid byte '"' in Cookie.Value; dropping invalid bytes
	// doesn't like the levelInstances cookie; maybe encode to %22?
	for _, c := range cooks {
		req.AddCookie(c)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return infos, err
	}

	r, err := util.ProcessResponse(res)

	if err := json.Unmarshal(r.Body, &infos); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return infos, err
}

func GetLevels(cooks []*http.Cookie) (LevelInfo, error) {
	var infos LevelInfo

	u, err := url.Parse("https://www.stockfighter.io/ui/levels")
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return infos, err
	}

	for _, c := range cooks {
		req.AddCookie(c)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return infos, err
	}

	r, err := util.ProcessResponse(res)
	if err = json.Unmarshal(r.Body, &infos); err != nil {
		fmt.Println(err)
		return infos, err
	}
	return infos, err
}
