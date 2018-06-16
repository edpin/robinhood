package robinhood

import "strings"

type Instrument string

func (i Instrument) GetID() string {
	fields := strings.Split(string(i), "/")
	if len(fields) < 2 {
		return ""
	}
	return fields[len(fields)-2]
}
