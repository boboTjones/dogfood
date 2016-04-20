//cab is for chock-a-block

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"sf"
	"util"
)

var (
	APIKEY                 string
	Res                    *util.Response
	Order                  *sf.Order
	Action, Method, Path   string
	Venue, Account, Symbol string
	Price, Qty, orderid    int
	err                    error
	data                   []byte
)

func init() {
	flag.StringVar(&Action, "do", Action, "buy, sell, check whatever")
	flag.StringVar(&Venue, "v", Venue, "Exchange venue name")
	flag.StringVar(&Account, "a", Account, "Account Number")
	flag.StringVar(&Symbol, "s", Symbol, "Symbol")
	flag.IntVar(&Price, "p", Price, "Price")
	flag.IntVar(&Qty, "q", Qty, "Quantity")
	flag.IntVar(&orderid, "i", orderid, "id of the order you want to modify")
}

func main() {

	APIKEY = os.Getenv("SF_APIKEY")
	if APIKEY == "" {
		fmt.Println("export SF_APIKEY=")
		os.Exit(2)
	}

	flag.Parse()

	t, err := url.Parse("https://api.stockfighter.io/")

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	switch Action {
	case "heartbeat":
		Method, Path = sf.Heartbeat()
	case "list":
		Method, Path = sf.GetStocks(Venue)
	case "quote":
		q, _ := sf.GetQuote(t, Venue, Symbol, APIKEY)
		fmt.Println(q)
	case "orders":
		Method, Path = sf.GetOrdersForAcct(Venue, Account)
	case "mine":
		Method, Path = sf.GetOrdersForAcctForSymbol(Venue, Account, Symbol)
	case "limit":
		Method, Path = sf.MakeOrder(Venue, Symbol)
		Order = sf.NewOrder(Account, Venue, Symbol, "buy", "limit", Price, Qty)
		data, err = json.Marshal(Order)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	case "market":
		Method, Path = sf.MakeOrder(Venue, Symbol)
		Order = sf.NewOrder(Account, Venue, Symbol, "buy", "market", Price, Qty)
		data, err = json.Marshal(Order)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	case "delete":
		Method, Path = sf.DelOrder(Venue, Symbol, orderid)
	case "restart":
		r, err := sf.RestartLastLevel()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		// for now, I'm going to assume each restart provides the same info.
		fmt.Printf("Restarted %g.\nVenues\t%s\nAccount\t%s\nTickers\t%s\n", r.InstanceId, r.Venues[0], r.Account, r.Tickers[0])
		os.Exit(1)
	default:
		fmt.Println("Gimme something to work with, ok?")
		os.Exit(2)
	}

	switch Method {
	case "GET":
		t.Path = Path
		Res, err = util.Get(t, APIKEY)
	case "POST":
		t.Path = Path
		Res, err = util.Post(t, data, APIKEY)
	case "DELETE":
		t.Path = Path
		Res, err = util.Del(t, APIKEY)
	default:
		fmt.Println("Method not supplied")
		os.Exit(2)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	fmt.Println(string(Res.Body))
}
