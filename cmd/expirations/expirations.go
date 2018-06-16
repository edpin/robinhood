// expirations returns all expirations for options of a given underlying symbol.
package main

import (
	"flag"
	"fmt"

	rh "github.com/edpin/robinhood"
)

var symbol = flag.String("symbol", "", "Symbol of the underlying security")

func main() {
	flag.Parse()

	if flag.NFlag() != 1 {
		fmt.Printf(`
Usage:
  expirations --symbol=F
`)
		return
	}
	client := &rh.Client{}

	exp, err := client.Expirations(*symbol)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Expirations: %+v\n", exp)
}
