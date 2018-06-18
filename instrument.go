package robinhood

import "strings"

// Instrument represents a Robinhood resource.
type Instrument string

// GetID returns the ID part of the Instrument.
func (i Instrument) GetID() string {
	fields := strings.Split(string(i), "/")
	if len(fields) < 2 {
		return ""
	}
	return fields[len(fields)-2]
}
