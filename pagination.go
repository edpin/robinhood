package robinhood

import (
	"bytes"
	"encoding/json"
)

// Deals with pagination transparently.

type paginatedResults struct {
	Results json.RawMessage `json:"results"`
	Next    string          `json:"next"`
}

func (c *Client) paginatedGet(endpoint string) ([]byte, error) {
	var buf bytes.Buffer
	for {
		resp, err := c.get(endpoint)
		if err != nil {
			return nil, err
		}
		var res paginatedResults
		err = json.Unmarshal(resp, &res)
		if err != nil {
			return nil, err
		}
		buf.Write([]byte(res.Results))
		if res.Next == "" {
			return buf.Bytes(), nil
		}
    endpoint = res.Next
	}
	// NOT REACHED
}
