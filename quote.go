package robinhood

// This file issues requests for stock price quotes.

import (
	"encoding/json"
	"strings"
)

// Quote is the real-time price quote for a security.
type Quote struct {
	Symbol string
	Ask    float64
	Bid    float64
}

// Quote returns a slice of quotes for the requested security symbols. Does
// not work on option symbols.
func (c *Client) Quote(symbol []string) ([]Quote, error) {
	quotes, err := c.quote(symbol)
	if err != nil {
		return nil, err
	}
	// Convert prices from string to floats.
	var qts []Quote
	for _, q := range quotes {
		bid, err := parseFloat64(q.Bid, nil)
		ask, err := parseFloat64(q.Ask, err)
		if err != nil {
			return nil, err
		}
		qts = append(qts, Quote{
			Symbol: q.Symbol,
			Ask:    ask,
			Bid:    bid,
		})
	}
	return qts, nil
}

type quote struct {
	Ask        string     `json:"ask_price"`
	Bid        string     `json:"bid_price"`
	Symbol     string     `json:"symbol"`
	Instrument Instrument `json:"instrument"`
}

func (c *Client) quote(symbol []string) ([]quote, error) {
	if len(symbol) == 0 {
		return nil, nil
	}
	var resp []byte
	var err error
	var quotes []quote
	var q quote
	if len(symbol) == 1 {
		resp, err = c.get(quotesURI + symbol[0] + "/")
	} else {
		resp, err = c.get(quotesURI + "?symbols=" + strings.Join(symbol, ","))
	}
	if err != nil {
		return nil, err
	}
	if len(symbol) == 1 {
		err = json.Unmarshal(resp, &q)
	} else {
		var results map[string][]quote
		err = json.Unmarshal(resp, &results)
		quotes = results["results"]
	}
	if err != nil {
		return nil, err
	}
	if len(symbol) == 1 {
		quotes = append(quotes, q)
	}
	return quotes, nil
}
