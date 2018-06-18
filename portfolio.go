package robinhood

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

// This file deals with portfolio information for a given account.

// Position identifies a position in a portfolio.
type Position struct {
	Symbol   string
	Name     string
	BuyPrice float64
	Quantity float64
  // TODO: add other fields.
}

// Portfolio returns a slice of Position a user has in their account.
func (c *Client) Portfolio() ([]Position, error) {
	var positions []Position
	pos, err := c.portfolio()
	if err != nil {
		return nil, err
	}
	for _, p := range pos {
		req, err := http.NewRequest("GET", p.URL, nil)
		if err != nil {
			return nil, err
		}
		resp, err := c.doReqWithAuth(req)
		if err != nil {
			log.Printf("Error fetching details for position %v: %v", p, err)
			continue
		}
		var detail detailedPosition
		err = json.Unmarshal(resp, &detail)
		if err != nil {
			log.Printf("Error unmarshalling details for position %v: %v", p, err)
			continue
		}
		buyPrice, err := parseFloat64(p.BuyPrice, nil)
		quantity, err := parseFloat64(p.Quantity, err)
		if err != nil {
			log.Printf("Error parsing float: %v", err)
			continue
		}
		positions = append(positions, Position{
			Symbol:   detail.Symbol,
			Name:     detail.Name,
			BuyPrice: buyPrice,
			Quantity: quantity,
		})
	}
	return positions, nil
}

type position struct {
	BuyPrice string `json:"average_buy_price"`
	URL      string `json:"instrument"`
	Quantity string `json:"quantity"`
}

func (c *Client) portfolio() ([]position, error) {
	parms := url.Values{}
	parms.Set("nonzero", "true")
	resp, err := c.paginatedGet(accountsURI + c.AccountID + "/" + positionsURI + "?" + parms.Encode())
	if err != nil {
		return nil, err
	}
	var positions []position
	err = json.Unmarshal(resp, &positions)
	if err != nil {
		return nil, err
	}
	return positions, nil
}

type detailedPosition struct {
	Symbol string `json:"symbol"`
	Name   string `json:"simple_name"`
}
