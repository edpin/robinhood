package robinhood

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Client is the Robinhood API client. It supports a single account. For users
// with multiple accounts, create a new Client for each account.
type Client struct {
	// AccountID is the account number this client will use. It is required.
	AccountID string

	// Token is the access token used for authentication. If it's not present,
	// then Username and Password must be present.
	Token string

	// Username is the user's username with Robinhood. It may be blank.
	Username string

	// Password is the user's password. It may be blank.
	Password string
}

type token struct {
	Token string `json:"token"`
}

// GetToken gets a new token, based on this client's Username and Password.
// It implicitly saves the new token.
func (c *Client) GetToken() error {
	if c.Username == "" || c.Password == "" {
		return fmt.Errorf("invalid username or password; neither can be blank")
	}

	form := url.Values{}
	form.Add("username", c.Username)
	form.Add("password", c.Password)

	resp, err := post(tokenURL, form.Encode(), "")
	if err != nil {
		return err
	}
	var tok token
	err = json.Unmarshal(resp, &tok)
	if err != nil {
		return err
	}
	if tok.Token == "" { // TODO: add other checks here, maybe length of token.
		return fmt.Errorf("invalid token returned: %v", tok.Token)
	}
	c.Token = tok.Token
	return nil
}

type Account struct {
  AccountNumber string `json:"account_number"`
}

type accounts struct {
  Accounts []Account `json:"results"`
  Next string `json:"next"`
}

// GetAccounts returns the list of all account numbers associated with a user.
// Client must be authenticated (i.e. a Token must be supplied).
func (c *Client) GetAccounts() ([]Account, error) {
  resp, err := get(accountsURL, c.Token)
  if err != nil {
    return nil, err
  }
  var accs accounts
  err = json.Unmarshal(resp, &accs)
  if err != nil {
    return nil, err
  }
  return accs.Accounts, nil
}
