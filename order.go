package robinhood

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
)

// This file deals with placing and canceling orders.

// Duration is the duration of a buy/sell order.
type Duration int

const (
	// Day means the order is valid for today only.
	Day Duration = iota

	// GTC means Good-Till-Cancelled.
	GTC
)

// OrderType describes whether the order is market, stop, limit, etc.
type OrderType int

// See description for OrderType.
const (
	Market OrderType = iota
	Limit
	Stop
	StopLimit
)

// Side represents the side of the order, either a buy or a sell and for short
// orders, the opening or closing side.
type Side int

// See description for Side.
const (
	Buy Side = iota
	Sell
	BuyToOpen
	BuyToClose
	SellToOpen
	SellToClose
)

// Order describes a buy or sell order.
type Order struct {
	Symbol    string
	Quantity  int64
	Duration  Duration
	Type      OrderType
	Side      Side
	Price     float64
	StopPrice float64 // only present for STOP or STOP_LIMIT orders.
}

// Order creates a new trade order for this client's account.
func (c *Client) Order(o Order) error {
	// Error checking
	if c.AccountID == "" {
		return fmt.Errorf("no account id provided in client")
	}
	if o.StopPrice < 0 {
		return fmt.Errorf("stop price must never be negative")
	}
	if o.Price <= 0.0001 {
		return fmt.Errorf("price must never be zero or negative")
	}
	// Find the instrument.
	quotes, err := c.quote([]string{o.Symbol})
	if err != nil {
		return err
	}
	if len(quotes) != 1 {
		return fmt.Errorf("invalid quote returned for symbol %q: %v", o.Symbol, err)
	}
	// Fetch account URL. This could be assembled from the appropriate URI pieces,
	// but this way is safer against trivial endpoint changes.
	accs, err := c.GetAccounts()
	if err != nil {
		return err
	}
	accountURL := ""
	for _, a := range accs {
		if a.AccountNumber == c.AccountID {
			accountURL = a.URL
		}
	}
	if accountURL == "" {
		return fmt.Errorf("invalid account number %s", c.AccountID)
	}
	instrument := quotes[0].Instrument
	oType := "market"
	if o.Type != Market {
		oType = "limit"
	}
	form := url.Values{}
	form.Add("account", accountURL)
	form.Add("instrument", string(instrument))
	form.Add("symbol", o.Symbol)
	form.Add("type", oType)
	form.Add("time_in_force", o.Duration.String())
	form.Add("trigger", "immediate")
	form.Add("price", fmt.Sprintf("%2.2f", o.Price))
	form.Add("quantity", fmt.Sprintf("%d", o.Quantity))
	form.Add("side", o.Side.String())
	if o.StopPrice != 0 {
		form.Set("trigger", "stop")
		form.Add("stop_price", fmt.Sprintf("%2.2f", o.StopPrice))
	}
	log.Printf("Posting order: %s", form.Encode())
	resp, err := c.post(ordersURI, form.Encode())
	if err != nil {
		return err
	}
	var status orderStatus
	err = json.Unmarshal(resp, &status)
	if err != nil {
		return err
	}
	// TODO: interpret order status. Or get rid of it?
	return nil
}

// String implements Stringer.
func (d Duration) String() string {
	switch d {
	case Day:
		return "gfd"
	case GTC:
		return "gtc"
	}
	return "(invalid duration)"
}

// String implements Stringer.
func (s Side) String() string {
	switch s {
	case Sell:
		return "sell"
	case Buy:
		return "buy"
	case BuyToOpen:
		return "buy_to_open"
	case BuyToClose:
		return "buy_to_close"
	case SellToOpen:
		return "sell_to_open"
	case SellToClose:
		return "sell_to_close"
	default:
		return "invalid side"
	}
}

// String implements Stringer.
func (t OrderType) String() string {
	switch t {
	case Market:
		return "market"
	case Limit:
		return "limit"
	case Stop:
		return "stop"
	case StopLimit:
		return "stop_limit"
	default:
		return "invalid order type"
	}
}

/*
cancel	URL	If this is not null, you can POST to this URL to cancel the order
id	String	Internal id of this order
reject_reason	String
state	String	queued, unconfirmed, confirmed, partially_filled, filled, rejected, canceled, or failed
*/
// TODO: currently not used. Remove?
type orderStatus struct {
	ID           string `json:"id"`
	Cancel       string `json:"cancel"`
	RejectReason string `json:"reject_reason"`
	State        string `json:"state"` // queued, unconfirmed, confirmed, partially_filled, filled, rejected, canceled, or failed
}
