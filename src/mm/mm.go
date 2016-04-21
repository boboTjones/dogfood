package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sf"
)

var (
	li                  *sf.LevelInfo
	APIKEY              string
	v, s, a             string
	action              string
	price, qty, orderId float64
	user, pass          string
)

func setup() *url.URL {
	APIKEY = os.Getenv("SF_APIKEY")
	if APIKEY == "" {
		fmt.Println("export SF_APIKEY=")
		os.Exit(2)
	}

	li, err := sf.GetLevels(APIKEY)
	if !li.OK {
		restart, err := sf.RestartLastLevel()
		if err != nil {
			fmt.Println("Restart failed ", err)
		}
		fmt.Println(restart)
		os.Exit(2)
	}
	li, err = sf.GetLevels(APIKEY)
	v = li.Venues[0]
	s = li.Tickers[0]
	a = li.Account

	t, err := url.Parse("https://api.stockfighter.io/")

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	return t
}

func init() {
	flag.StringVar(&action, "do", action, "buy, sell, quote")
	flag.Float64Var(&price, "c", price, "price")
	flag.Float64Var(&qty, "q", price, "qty")
	flag.Float64Var(&orderId, "i", orderId, "order of the id you wish to cancel")
	flag.StringVar(&user, "u", user, "username")
	flag.StringVar(&pass, "p", pass, "password")
}

func main() {
	flag.Parse()

	//ok := sf.Vbeat(t, v, APIKEY)
	switch action {
	case "login":
		if user == "" || pass == "" {
			fmt.Println("Invalid username or password.")
			os.Exit(2)
		}
		cookies, err := sf.Login(user, pass)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(cookies)
	case "buy":
		t := setup()
		q, _ := sf.GetQuote(t, v, s, APIKEY)
		if price == 0 {
			price = q.Last - 300.00
		}
		if qty == 0 {
			fmt.Println("How many?")
			os.Exit(2)
		}
		o := sf.NewOrder(a, v, s, "buy", "limit", price, qty)
		data, err := json.Marshal(o)
		foo, err := sf.MakeOrder(t, data, v, s, APIKEY)
		fmt.Println(string(foo.Body))
		fmt.Println(err)
	case "sell":
		t := setup()
		if qty == 0 {
			fmt.Println("How many?")
			os.Exit(2)
		}
		o := sf.NewOrder(a, v, s, "sell", "limit", price, qty)
		data, err := json.Marshal(o)
		foo, err := sf.MakeOrder(t, data, v, s, APIKEY)
		fmt.Println(string(foo.Body))
		fmt.Println(err)
	case "quote":
		t := setup()
		_, _ = sf.GetQuote(t, v, s, APIKEY)
	case "info":
		t := setup()
		i, _ := sf.StockInfo(t, v, s, APIKEY)
		fmt.Println(i)
	case "show":
		t := setup()
		trades, _ := sf.Mine(t, v, a, s, APIKEY)
		fmt.Println(trades.Orders)
	case "cancel":
		t := setup()
		if orderId == 0 {
			fmt.Println("Must provide order id")
			os.Exit(2)
		}
		sf.Cancel(t, v, s, APIKEY, orderId)
	case "restart":
		restart, err := sf.RestartLastLevel()
		fmt.Println(err)
		fmt.Println(restart)
	default:
		fmt.Println("BZZT TRY AGAIN")
	}

}
