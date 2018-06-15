// portfolio retrieves the user's portfolio.
package main

import (
	"flag"
	"fmt"

	rh "github.com/edpin/robinhood"
)

var (
	account = flag.String("account", "", "User's account with Robinhood")
	token   = flag.String("token", "", "User's access token with Robinhood")
)

func main() {
	flag.Parse()

	if flag.NFlag() != 2 {
		fmt.Printf(`
Usage:
  portfolio --account=<account_id> --token=<access_token>
`)
		return
	}

	client := &rh.Client{
		AccountID: *account,
		Token:     *token,
	}

	port, err := client.Portfolio()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Got: %v\n", port)
}
