package robinhood

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// Deals with pagination transparently.

type paginatedResults struct {
	Results json.RawMessage `json:"results"`
	Next    string          `json:"next"`
}

func (c *Client) paginatedGet(endpoint string) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("[")
	endpoint = apiURL + endpoint
	for {
		req, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.doReqWithAuth(req)
		if err != nil {
			return nil, err
		}
		var res paginatedResults
		err = json.Unmarshal(resp, &res)
		if err != nil {
			return nil, err
		}
		// TODO: This is inefficient. Fix it!
		resStr := strings.Trim(string(res.Results), "[] ")
		buf.Write([]byte(resStr))
		if res.Next == "" {
			buf.WriteString("]")
			return buf.Bytes(), nil
		}
		if resStr != "" {
			buf.WriteString(",")
		}
		endpoint = res.Next
	}
	// NOT REACHED
}
