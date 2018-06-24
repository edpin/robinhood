# Robinhood Go API

Work In Progress.

This library is under development. It is already functional for the following:

- Initialize with user's credentials.
- Fetch portfolio and account information.
- Get real-time quotes.
- Get options chains.
- Enter simple stock orders.

TODO:

- Option orders.
- Multi-leg option orders.
- More testing.

## To start using:

```
import (
  "fmt"
  
  rh "github.com/edpin/robinhood"
)

func main() {
  client := &rh.Client{
		Username: "username",
		Password: "password",
  }

  err := client.GetToken()
  if err != nil {
		panic(err)
  }

	accs, err := client.GetAccounts()
	if err != nil {
		panic(err)
	}
	fmt.Println("Accounts:")
	for _, acc := range accs {
		fmt.Printf("Account: %v\n", acc.AccountNumber)
	}
}
```
See cmd for more examples.
