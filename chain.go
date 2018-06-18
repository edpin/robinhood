package robinhood

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"time"
)

// This file deals with retrieving and parsing option chains information.

// Chain represents a complete option with strike price, expiration, type
// (put/call) and an underlying symbol.
type Chain struct {
	Strike     float64
	Expiration time.Time
	Symbol     string
	Type       string // "put" or "call". TODO: use an enum?

	// Private fields
	id string
}

// Chains returns all chains (i.e. a complete option with strike price) for
// an option on the underlying symbol and an expiration date.
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

type chain struct {
	ID             string     `json:"id"`
	InstrumentID   Instrument `json:"instrument"`
	ChainID        string     `json:"chain_id"`
	URL            string     `json:"url"` // url to quote?
	StrikePrice    string     `json:"strike_price"`
	ExpirationDate string     `json:"expiration_date"`
	Type           string     `json:"type"` // "put" or "call"
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
