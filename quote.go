package robinhood

// Issues requests for price quotes.

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type quote struct {
	Ask        string     `json:"ask_price"`
	Bid        string     `json:"bid_price"`
	Symbol     string     `json:"symbol"`
	Instrument Instrument `json:"instrument"`
}

type Quote struct {
	Ask    float64
	Bid    float64
	Symbol string
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
			Ask:    ask,
			Bid:    bid,
			Symbol: q.Symbol,
		})
	}

	return qts, nil
}

func (c *Client) getInstrumentID(symbol string) (string, error) {
	quotes, err := c.quote([]string{symbol})
	if err != nil {
		return "", err
	}
	if len(quotes) != 1 {
		return "", fmt.Errorf("expected one quote for symbol %s, got %d", symbol, len(quotes))
	}
	return quotes[0].Instrument.GetID(), nil
}

type expirations struct {
	ID          string   `json:"id"`
	Symbol      string   `json:"symbol"`
	Expirations []string `json:"expiration_dates"`
}

const dateFormat = "2006-01-02"

func (c *Client) expirations(symbol string) (expirations, error) {
	var e0 expirations
	instrumentID, err := c.getInstrumentID(symbol)
	if err != nil {
		return e0, err
	}
	resp, err := c.paginatedGet(chainsURI + "?equity_instrument_ids=" + instrumentID)
	if err != nil {
		return e0, err
	}
	var exp []expirations
	err = json.Unmarshal(resp, &exp)
	if err != nil {
		return e0, err
	}
	// Search for the correct symbol here as there are often things like 1SPY,
	// 2SPY, SPY and 1GOOG, 2GOOG, GOOG, etc.
	for _, e := range exp {
		if e.Symbol == symbol {
			return e, nil
		}
	}
	return e0, fmt.Errorf("not found instrumentID %q for symbol %q", instrumentID, symbol)
}

// Expirations returns all expiration dates for options for the underlying symbol.
func (c *Client) Expirations(symbol string) ([]time.Time, error) {
	exp, err := c.expirations(symbol)
	if err != nil {
		return nil, err
	}
	var exps []time.Time
	for _, date := range exp.Expirations {
		exp, err := time.Parse(dateFormat, date)
		if err != nil {
			log.Printf("Error converting expiration %s: %s", date, err)
			continue
		}
		exps = append(exps, exp)
	}
	return exps, nil
}

type chain struct {
	ID             string     `json:"id"`
	InstrumentID   Instrument `json:"instrument"`
	ChainID        string     `json:"chain_id"`
	URL            string     `json:"url"` // url to quote?
	StrikePrice    string     `json:"strike_price"`
	ExpirationDate string     `json:"expiration_date"`
	Type           string     `json:"type"` // "put" or "call"
}

type Chain struct {
	Strike     float64
	Expiration time.Time
	Symbol     string
	Type       string // "put" or "call". TODO: use an enum?

	// Private fields
	id string
}

func (c *Client) chains(symbol string, expiration time.Time) ([]chain, error) {
	exp, err := c.expirations(symbol)
	if err != nil {
		return nil, err
	}

	// Fetch the strikes for this option ID at a given expiration date.
	parms := url.Values{}
	parms.Set("chain_id", exp.ID)
	parms.Set("expiration_dates", expiration.Format(dateFormat))
	parms.Set("state", "active")
	parms.Set("tradability", "tradable")

	resp, err := c.paginatedGet(optionsURI + "?" + parms.Encode())
	if err != nil {
		return nil, err
	}
	var chains []chain
	err = json.Unmarshal(resp, &chains)
	if err != nil {
		return nil, err
	}
	return chains, nil
}

func (c *Client) Chains(symbol string, expiration time.Time) ([]Chain, error) {
	chains, err := c.chains(symbol, expiration)
	if err != nil {
		return nil, err
	}
	// Convert internal format to external format.
	var Chains []Chain
	for _, c := range chains {
		strike, err := strconv.ParseFloat(c.StrikePrice, 64)
		if err != nil {
			log.Printf("error converting to float %q: %v", c.StrikePrice, err)
			continue
		}
		exp, err := time.Parse(dateFormat, c.ExpirationDate)
		if err != nil {
			log.Printf("error parsing expiration date %q: %v", c.ExpirationDate, err)
			continue
		}
		Chains = append(Chains, Chain{
			Symbol:     symbol,
			Type:       c.Type,
			Strike:     strike,
			Expiration: exp,
			id:         c.ID,
		})
	}
	return Chains, nil
}

type option struct {
	Ask         string `json:"ask_price"`
	Bid         string `json:"bid_price"`
	Last        string `json:"last_trade_price"`
	MarketPrice string `json:"adjusted_mark_price"`
	IV          string `json:"implied_volatility"`
	// Other fields available...
}

type Option struct {
	Symbol      string // the underlying symbol.
	Strike      float64
	Expiration  time.Time
	Type        string // "put" or "call". TODO: use an enum?
	Bid         float64
	Ask         float64
	Last        float64
	MarketPrice float64 // Close to midpoint, but not quite.
	IV          float64
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
	log.Printf("Got raw option: %s", resp)
	var o option
	err = json.Unmarshal(resp, &o)
	if err != nil {
		return o0, err
	}
	bid, err := parseFloat64(o.Bid, nil)
	ask, err := parseFloat64(o.Ask, err)
	last, err := parseFloat64(o.Last, err)
	marketPrice, err := parseFloat64(o.MarketPrice, err)
	//iv, err := parseFloat64(o.IV, err)
	option := Option{
		Symbol:      chain.Symbol,
		Strike:      chain.Strike,
		Expiration:  chain.Expiration,
		Type:        chain.Type,
		Bid:         bid,
		Ask:         ask,
		Last:        last,
		MarketPrice: marketPrice,
		//IV:          iv,
	}
	return option, err
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
