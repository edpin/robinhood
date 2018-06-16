// chains returns option chains of a given underlying symbol and expiration.
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
)

func main() {
	flag.Parse()

	if flag.NFlag() != 2 {
		fmt.Printf(`
Usage:
  chains --symbol=SPY --expiration=2018-06-22
`)
		return
	}
	client := &rh.Client{}

	exp, err := time.Parse("2006-01-02", *expiration)
	if err != nil {
		panic(err)
	}
	c, err := client.Chains(*symbol, exp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Chains: %+v\n", c)
}
