// Package robinhood is a client of Robinhood's trading and screening API.
package robinhood

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	apiURL       = "https://api.robinhood.com/"
	tokenURI     = "api-token-auth/"
	accountsURI  = "accounts/"
	positionsURI = "positions/"
)

// get performs an HTTP get request on 'endpoint' using 'token' for authentication.
func get(endpoint, token string) ([]byte, error) {
	req, err := http.NewRequest("GET", apiURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	return doReq(req, token)
}

// post performs an HTTP post of 'data' to 'endpoint' using 'token' for authentication.
func post(endpoint string, data, token string) ([]byte, error) {
	buf := strings.NewReader(data)
	req, err := http.NewRequest("POST", apiURL+endpoint, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// Note: Content-Length is set by NewRequest.
	return doReq(req, token)
}

func doReq(req *http.Request, token string) ([]byte, error) {
	req.Header.Add("Accept", "application/json")
	if token != "" {
		req.Header.Add("Authorization", "Token "+token)
	}
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
