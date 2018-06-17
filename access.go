// Package robinhood is a client of Robinhood's trading and screening API.
package robinhood

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiURL           = "https://api.robinhood.com/"
	tokenURI         = "api-token-auth/"
	accountsURI      = "accounts/"
	positionsURI     = "positions/"
	quotesURI        = "quotes/"
	chainsURI        = "options/chains/"      // ?equity_instrument_ids=
	optionsURI       = "options/instruments/" //?chain_id={_chainid}&expiration_dates={_dates}&state=active&tradability=tradable
	marketOptionsURI = "marketdata/options/"  //{_optionid}/
	oAuthUpgradeURI  = "oauth2/migrate_token/"
)

// get performs an HTTP get request on 'endpoint'..
func (c *Client) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.doReqWithAuth(req)
}

// post performs an HTTP post of 'data' to 'endpoint'.
func (c *Client) post(endpoint string, data string) ([]byte, error) {
	buf := strings.NewReader(data)
	req, err := http.NewRequest("POST", apiURL+endpoint, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// Note: Content-Length is set by NewRequest.
	return c.doReqWithAuth(req)
}

func (c *Client) doReqWithAuth(req *http.Request) ([]byte, error) {
	_, hasAuthorizationHeader := req.Header["Authorization"]
	if c.Token != "" && !hasAuthorizationHeader {
		req.Header.Add("Authorization", "Token "+c.Token)
	}
	return c.doReq(req)
}

func (c *Client) doReq(req *http.Request) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Request failed: %s (%d)", resp.Status, resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) doReqWithBearerToken(req *http.Request) ([]byte, error) {
	err := c.EnsureBearerToken()
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.BearerToken)
	return c.doReq(req)
}
