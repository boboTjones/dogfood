package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	data map[string]interface{}
)

type Response struct {
	Path    string
	Method  string
	Code    int
	Body    []byte
	Err     error
	Cookies map[string]string
}

type Order struct {
	Account   string `json:"account"`
	Venue     string `json:"venue"`
	Stock     string `json:"stock"`
	Price     int    `json:"price"`
	Qty       int    `json:"qty"`
	Direction string `json:"direction"`
	OrderType string `json:"orderType"`
}

func SlurpFromURL(t string) []byte {
	r, err := http.Get(t)
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Something bad happened.")
	}
	return b
}

func SlurpFromFile(filePath string) []byte {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Something bad happened: %v", err)
	}
	return data
}

func Post(u *url.URL, data []byte, key string) (*Response, error) {
	req, err := http.NewRequest("POST", u.String(),
		bytes.NewReader(data))
	if err != nil {
		return &Response{Path: u.Path, Method: "POST", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)
	fmt.Println(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "POST", Err: err}, err
	}

	r, err := processResponse(res)
	r.Path = u.Path
	r.Method = "POST"
	r.Err = err
	return r, err
}

func Get(u *url.URL, key string) (*Response, error) {
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return &Response{Path: u.Path, Method: "GET", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)
	fmt.Println(req)

	client := http.Client{Timeout: (2999 * time.Millisecond)}

	res, err := client.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "GET", Err: err}, err
	}

	r, err := processResponse(res)
	r.Path = u.Path
	r.Method = "GET"
	r.Err = err
	return r, err
}

func Put(u *url.URL, content []byte, key string) (*Response, error) {
	req, err := http.NewRequest("PUT", u.String(), bytes.NewReader(content))
	if err != nil {
		return &Response{Path: u.Path, Method: "PUT", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "PUT", Err: err}, err
	}

	r, err := processResponse(res)
	r.Path = u.Path
	r.Method = "PUT"
	r.Err = err
	return r, err
}

func Del(u *url.URL, key string) (*Response, error) {
	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return &Response{Path: u.Path, Method: "DELETE", Err: err}, err
	}

	req.Header.Set("X-Starfighter-Authorization", key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &Response{Path: u.Path, Method: "DELETE", Err: err}, err
	}

	r, err := processResponse(res)
	r.Path = u.Path
	r.Method = "DELETE"
	r.Err = err
	return r, err
}

func processResponse(res *http.Response) (r *Response, err error) {
	ret := &Response{
		Code: res.StatusCode,
	}

	ret.Body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func Heartbeat() (string, string) {
	//get /venues/:venue/heartbeat
	return "GET", "/ob/api/heartbeat"
}

func Vbeat(v string) (string, string) {
	//get /venues/:venue/heartbeat
	return "GET", "/ob/api/venues/" + v + "/heartbeat"
}

func GetStocks(v string) (string, string) {
	//get /venues/:venue/stocks
	return "GET", "/ob/api/venues/" + v + "/stocks"
}

func GetStock(v, s string) (string, string) {
	//get /venues/:venue/stocks/:stock
	return "GET", "/ob/api/venues/" + v + "/stocks/" + s
}

func MakeOrder(v, s string) (string, string) {
	//post /venues/:venue/stocks/:stock/orders
	return "POST", "/ob/api/venues/" + v + "/stocks/" + s + "/orders"
}

func GetQuote(v, s string) (string, string) {
	//get /venues/:venue/stocks/:stock/quote
	return "GET", "/ob/api/venues/" + v + "/stocks/" + s + "/quote"
}

func GetOrder(v, s, id string) (string, string) {
	//get /venues/:venue/stocks/:stock/orders/:id
	return "GET", "/ob/api/venues/" + v + "/stocks/" + s + "/orders" + id

}

func DelOrder(v, s, id string) (string, string) {
	//delete /venues/:venue/stocks/:stock/orders/:order
	return "DELETE", "/ob/api/venues/" + v + "/stocks/" + s + "/orders" + id

}

func GetOrdersForAcct(v, a string) (string, string) {
	//get /venues/:venue/accounts/:account/orders
	return "GET", "/ob/api/venues/" + v + "/accounts/" + a + "/orders"

}

func GetOrdersForAcctForSymbol(v, a, s string) (string, string) {
	//get /venues/:venue/accounts/:account/stocks/:stock/orders
	return "GET", "/ob/api/venues/" + v + "/accounts/" + a + "/stocks" + s + "/orders"

}

func NewOrder(a, v, s, d, ot string, p, qty int) *Order {
	return &Order{
		Account:   a,
		Venue:     v,
		Stock:     s,
		Direction: d,
		OrderType: ot,
		Price:     p,
		Qty:       qty,
	}
}

func RestartLastLevel() (map[string]interface{}, error) {
	cooks := make(map[string]string)
	var levels map[string]interface{}
	var level string
	var infos map[string]interface{}

	x := SlurpFromFile("/Users/erin/cookies.txt")
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
	for k, _ := range levels {
		level = k
	}

	u, _ := url.Parse("https://www.stockfighter.io/gm/levels/" + level)
	req, err := http.NewRequest("POST", u.String(), nil)
	if err != nil {
		return infos, err
	}
	for k, v := range cooks {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return infos, err
	}

	r, err := processResponse(res)

	if err := json.Unmarshal(r.Body, &infos); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return infos, err
}
