// get_token is a simple utility to get a user's token, given their username and password.
package main

import (
	"flag"
	"fmt"

	rh "github.com/edpin/robinhood"
)

var (
	username = flag.String("username", "", "Username with Robinhood")
	password = flag.String("password", "", "Password with Robinhood")
)

func main() {
	flag.Parse()

	if flag.NFlag() != 2 {
		fmt.Printf(`
Usage:
  get_token --username=your_user_name --password=your_password
`)
		return
	}

	client := &rh.Client{
		Username: *username,
		Password: *password,
	}

	fmt.Println("Fetching token for user ", *username)

	err := client.GetToken()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %s\n", client.Token)
}
