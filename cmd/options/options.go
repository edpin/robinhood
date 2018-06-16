// options returns live quota for the specified option.
package main

import (
	"flag"
	"fmt"
	"time"

	rh "github.com/edpin/robinhood"
)

var (
	symbol     = flag.String("symbol", "", "Symbol of the underlying security")
	expiration = flag.String("expiration", "", "Expiration date (YYYY-MM-DD)")
	strike     = flag.Float64("strike", 0.0, "Strike price of option")
	account    = flag.String("account", "", "User's account with Robinhood")
	token      = flag.String("token", "", "User's access token with Robinhood")
	optType    = flag.String("type", "put", "<put|call>")
)

func main() {
	flag.Parse()

	if flag.NFlag() < 5 {
		fmt.Printf(`
Usage:
  options --symbol=SPY --expiration=2018-06-22 --strike=290 --account=<account> \
          --token=<auth_token>  --type=call
`)
		return
	}
	client := &rh.Client{
		AccountID: *account,
		Token:     *token,
	}
	exp, err := time.Parse("2006-01-02", *expiration)
	if err != nil {
		panic(err)
	}
	opt, err := client.Option(*symbol, exp, *strike, *optType)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Option: %+v\n", opt)
}
