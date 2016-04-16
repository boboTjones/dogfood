//cab is for chock-a-block

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"util"
)

var (
	APIKEY                 string
	Res                    *util.Response
	Order                  *util.Order
	Action, Method, Path   string
	Venue, Account, Symbol string
	Price, Qty             int
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
		Method, Path = util.Heartbeat()
	case "list":
		Method, Path = util.GetStocks(Venue)
	case "quote":
		Method, Path = util.GetQuote(Venue, Symbol)
	case "orders":
		Method, Path = util.GetOrdersForAcct(Venue, Account)
	case "mine":
		Method, Path = util.GetOrdersForAcctForSymbol(Venue, Account, Symbol)
	case "buy":
		Method, Path = util.MakeOrder(Venue, Symbol)
		Order = util.NewOrder(Account, Venue, Symbol, "buy", "limit", Price, Qty)
		data, err = json.Marshal(Order)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	case "restart":
		r, err := util.RestartLastLevel()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		// for now, I'm going to assume each restart provides the same info.
		fmt.Printf("Restarted %g.\nVenues\t%s\nAccount\t%s\nTickers\t%s\n", r["instanceId"], r["venues"], r["account"], r["tickers"])
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
