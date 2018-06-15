package robinhood

// Issues requests for price quotes.

import (
	"encoding/json"
	"strconv"
	"strings"
)

type quote struct {
	Ask    string `json:"ask_price"`
	Bid    string `json:"bid_price"`
	Symbol string `json:"symbol"`
}

type Quote struct {
	Ask    float64
	Bid    float64
	Symbol string
}

func (c *Client) Quote(symbol []string) ([]Quote, error) {
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

	// Now convert prices from string to floats.
	var qts []Quote
	for _, q = range quotes {
		bid, err := strconv.ParseFloat(q.Bid, 64)
		if err != nil {
			return nil, err
		}
		ask, err := strconv.ParseFloat(q.Ask, 64)
		if err != nil {
			return nil, err
		}
		qts = append(qts, Quote{Ask: ask, Bid: bid, Symbol: q.Symbol})
	}

	return qts, nil
}
