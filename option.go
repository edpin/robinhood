package robinhood

import (
	"encoding/json"
	"net/http"
	"time"
)

// This file deals with retrieving and parsing Option quotes.

// Option represents a real-time quote for an option.
type Option struct {
	Symbol      string // the underlying symbol.
	Strike      float64
	Expiration  time.Time
	Type        string // "put" or "call". TODO: use an enum?
	Bid         float64
	Ask         float64
	Last        float64
	MarketPrice float64 // Close to midpoint, but not quite.
	Volume      int64
	IV          float64
	// TODO: add other fields such as AskSize, BidSize, greeks, etc.

	instrument Instrument // URL of instrument. Used for placing orders.
}

// Option returns a quote for an option chain.
func (c *Client) Option(chain Chain) (Option, error) {
	var o0 Option

	req, err := http.NewRequest("GET", apiURL+marketOptionsURI+chain.id+"/", nil)
	if err != nil {
		return o0, err
	}
	resp, err := c.doReqWithBearerToken(req)
	if err != nil {
		return o0, err
	}
	var o option
	err = json.Unmarshal(resp, &o)
	if err != nil {
		return o0, err
	}
	bid, err := parseFloat64(o.Bid, nil)
	ask, err := parseFloat64(o.Ask, err)
	last, err := parseOptionalFloat64(o.Last, err)
	marketPrice, err := parseFloat64(o.MarketPrice, err)
	iv, err := parseOptionalFloat64(o.IV, err)
	option := Option{
		Symbol:      chain.Symbol,
		Strike:      chain.Strike,
		Expiration:  chain.Expiration,
		Type:        chain.Type,
		Bid:         bid,
		Ask:         ask,
		Last:        last,
		Volume:      o.Volume,
		MarketPrice: marketPrice,
		IV:          iv,
		instrument:  Instrument(o.Instrument),
	}
	return option, err
}

type option struct {
	Ask         string `json:"ask_price"`
	Bid         string `json:"bid_price"`
	Last        string `json:"last_trade_price"`
	MarketPrice string `json:"adjusted_mark_price"`
	Volume      int64  `json:"volume"`
	IV          string `json:"implied_volatility"`
	Instrument  string `json:"instrument"`
	// Other fields available...
}
