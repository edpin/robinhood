package robinhood

// Issues requests for price quotes.

import (
	"encoding/json"
	"fmt"
	"log"
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
		bid, err := strconv.ParseFloat(q.Bid, 64)
		if err != nil {
			return nil, err
		}
		ask, err := strconv.ParseFloat(q.Ask, 64)
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
	Type       string // "put" or "call". TODO: use an enum?
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
			Type:       c.Type,
			Strike:     strike,
			Expiration: exp,
		})
	}
	return Chains, nil
}

type Option struct {
	Bid float64
	Ask float64
}

func (c *Client) Option(symbol string, expiration time.Time, strike float64, optionType string) (Option, error) {
	var o0 Option
	chains, err := c.chains(symbol, expiration)
	if err != nil {
		return o0, err
	}
	// Look for the right strike price.
	for _, chain := range chains {
		strikeChain, err := strconv.ParseFloat(chain.StrikePrice, 64)
		if err != nil {
			log.Printf("error converting to float %q: %v", chain.StrikePrice, err)
			continue
		}
		if floatEquals(strike, strikeChain) {
			return c.option(chain.ID)
		}
	}
	return o0, fmt.Errorf("error finding option quote for symbol %q, expiration %s, strike %2.2f", symbol, expiration.Format(dateFormat), strike)
}

func (c *Client) option(optionID string) (Option, error) {
	var o0 Option

	log.Printf("Going to fetch: %s", marketOptionsURI+optionID+"/")
	resp, err := c.get(marketOptionsURI + optionID + "/")
	if err != nil {
		return o0, err
	}
	log.Printf("Got raw option: %s", resp)
	return o0, nil
}
