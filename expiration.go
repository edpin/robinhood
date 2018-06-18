package robinhood

// This file deals with option expirations.

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

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
