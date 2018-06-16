package robinhood

import "testing"

func TestGetID(t *testing.T) {
	i := Instrument("https://api.robinhood.com/instruments/8f92e76f-1e0e-4478-8580-16a6ffcfaef5/")
	if got, expected := i.GetID(), "8f92e76f-1e0e-4478-8580-16a6ffcfaef5"; got != expected {
		t.Fatalf("got = %q, expected = %q", got, expected)
	}
}

func TestFail(t *testing.T) {
	i := Instrument("foobar")
	if i.GetID() != "" {
		t.Fatalf(`expected "", got %q`, i.GetID())
	}
}
