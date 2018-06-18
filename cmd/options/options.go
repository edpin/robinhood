// options returns live quota for the specified option.
package main

import (
	"flag"
	"fmt"
	"math"
	"time"

	rh "github.com/edpin/robinhood"
)

var (
	symbol     = flag.String("symbol", "", "Symbol of the underlying security")
	expiration = flag.String("expiration", "", "Expiration date (YYYY-MM-DD)")
	strike     = flag.Float64("strike", 0.0, "Strike price of option")
	token      = flag.String("token", "", "User's access token with Robinhood")
	optType    = flag.String("type", "put", "<put|call>")
)

func main() {
	flag.Parse()

	if flag.NFlag() < 4 {
		fmt.Printf(`
Usage:
  options --symbol=SPY --expiration=2018-06-22 --strike=290 \
          --token=<auth_token>  --type=call
`)
		return
	}
	client := &rh.Client{
		Token: *token,
	}
	exp, err := time.Parse("2006-01-02", *expiration)
	if err != nil {
		panic(err)
	}
	chains, err := client.Chains(*symbol, exp)
	if err != nil {
		panic(err)
	}
	// Look for the right strike price.
	idx := -1
	for i, chain := range chains {
		if floatEquals(*strike, chain.Strike) {
			idx = i
			break
		}
	}
	if idx < 0 {
		panic("no chain found")
	}

	opt, err := client.Option(chains[idx])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", opt)
}

const fEpsilon = 0.00001

func floatEquals(a, b float64) bool {
	return math.Abs(a-b) < fEpsilon
}
