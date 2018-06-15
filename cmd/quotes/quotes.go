// quotes retrieves quotes from the Robinhood API for a given symbol(s).
package main

import (
	"flag"
	"fmt"
	"strings"

	rh "github.com/edpin/robinhood"
)

var (
	symbol = flag.String("symbol", "", "Comma-separated list of symbols to retrieve")
)

func main() {
	flag.Parse()

	if flag.NFlag() != 1 {
		fmt.Printf(`
Usage:
  quotes --symbol=GOOG,MSFT,GM
`)
		return
	}

	client := &rh.Client{}

	symbs := strings.TrimRight(*symbol, ", ")
	symbols := strings.Split(symbs, ",")

	quotes, err := client.Quote(symbols)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got: %+v\n", quotes)
}
