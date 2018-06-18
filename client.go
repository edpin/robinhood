package robinhood

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Client is the Robinhood API client. It supports a single account. For users
// with multiple accounts, create a new Client for each account.
type Client struct {
	// AccountID is the account number this client will use. It is required for
	// all operations that operate directly on a user's account, such as calls to
	// get the portfolio and entering and cancelling orders.
	AccountID string

	// Token is the access token used for authentication. If it's not present,
	// then Username and Password must be present.
	Token string

	// Username is the user's username with Robinhood. It may be blank if the
	// client has a Token already.
	Username string

	// Password is the user's password. It may be blank if the client has a Token
	// already.
	Password string

	// BearerToken is present when a call to GetBearerToken is successful. It is
	// only necessary for real-time quotes.
	BearerToken string

	// BearerTokenExpiration is the wall clock time that the bearer token expires.
	BearerTokenExpiration time.Time
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

	resp, err := c.post(tokenURI, form.Encode())
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

// Account contains the user's AccountNumber.
type Account struct {
	AccountNumber string `json:"account_number"`
}

// GetAccounts returns the list of all account numbers associated with a user.
// Client must be authenticated (i.e. a Token must be supplied).
func (c *Client) GetAccounts() ([]Account, error) {
	resp, err := c.paginatedGet(accountsURI)
	if err != nil {
		return nil, err
	}
	var accs []Account
	err = json.Unmarshal(resp, &accs)
	if err != nil {
		return nil, err
	}
	return accs, nil
}

/*
   "token_type": "Bearer",
   "access_token": "9Lg%WiectYtobuiewceIVUnhjiBGLUIeytekLBGJKDHGfvhjkfkuggbusfhukewrygfubasd",
   "expires_in": 300,
   "refresh_token": "BKLtvuglkYUV67VIbtuiE5cyFVHwerCWRT",
   "scope": "internal"
*/
type oAuthToken struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"` // seconds
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

// GetBearerToken fetches the bearer token and stores it implicitly.
func (c *Client) GetBearerToken() error {
	resp, err := c.post(oAuthUpgradeURI, "")
	if err != nil {
		return err
	}
	var oauth oAuthToken
	err = json.Unmarshal(resp, &oauth)
	if err != nil {
		return err
	}
	if oauth.TokenType == "Bearer" && oauth.AccessToken != "" {
		c.BearerToken = oauth.AccessToken
		c.BearerTokenExpiration = time.Now().Add(time.Duration(oauth.ExpiresIn) * time.Second)
		return nil
	}
	return fmt.Errorf("no bearer token in reply: %s", resp)
}

// EnsureBearerToken ensures the client has a bearer token with at least another
// 30 seconds of time to live.
func (c *Client) EnsureBearerToken() error {
	// Do we still have 30 seconds left to use the token?
	if c.BearerTokenExpiration.After(time.Now().Add(30 * time.Second)) {
		return nil
	}
	return c.GetBearerToken()
}

// parseFloat64 parses the float and returns the prevErr if non null or the
// current error. Use it to chain several calls without having to check for
// errors until the end of the chain.
func parseFloat64(str string, prevErr error) (float64, error) {
	f, err := strconv.ParseFloat(str, 64)
	if prevErr != nil {
		return f, prevErr
	}
	return f, err
}

// parseOptionalFloat64 is the same as parseFloat64 but empty strings return 0
// instead of causing an error.
func parseOptionalFloat64(str string, prevErr error) (float64, error) {
	if str == "" {
		return 0, prevErr
	}
	return parseFloat64(str, prevErr)
}
