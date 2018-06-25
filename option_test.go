package robinhood

import (
	"strings"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var options = map[string]string{
	apiURL + quotesURI + "SPY/":                                                                                                          `{"ask_price":"274.5500","bid_price":"274.5000","symbol":"SPY","instrument":"https://api.robinhood.com/instruments/8f92e76f-1e0e-4478-8580-16a6ffcfaef5/"}`,
	apiURL + chainsURI + "?equity_instrument_ids=8f92e76f-1e0e-4478-8580-16a6ffcfaef5":                                                   `{"previous":null,"results":[{"can_open_position":true,"symbol":"SPY","trade_value_multiplier":"100.0000","underlying_instruments":[{"instrument":"https:\/\/api.robinhood.com\/instruments\/8f92e76f-1e0e-4478-8580-16a6ffcfaef5\/","id":"6c3bf803-ec29-41c1-b721-3471351fc61d","quantity":100}],"expiration_dates":["2018-06-27","2018-06-29","2018-07-02","2018-07-03","2018-07-06","2018-07-09","2018-07-11","2018-07-13","2018-07-16","2018-07-18","2018-07-20","2018-07-23","2018-07-25","2018-07-27","2018-07-30","2018-08-03","2018-08-17","2018-09-21","2018-09-28","2018-10-19","2018-12-21","2018-12-31","2019-01-18","2019-03-15","2019-03-29","2019-06-21","2019-09-20","2019-12-20","2020-01-17","2020-03-20","2020-06-19","2020-12-18"],"cash_component":null,"min_ticks":{"cutoff_price":"0.00","below_tick":"0.01","above_tick":"0.01"},"id":"c277b118-58d9-4060-8dc5-a3b5898955cb"},{"can_open_position":false,"symbol":"2SPY","trade_value_multiplier":"100.0000","underlying_instruments":[{"instrument":"https:\/\/api.robinhood.com\/instruments\/8f92e76f-1e0e-4478-8580-16a6ffcfaef5\/","id":"ca3c00f6-2477-485e-940b-90f173ada716","quantity":100}],"expiration_dates":[],"cash_component":null,"min_ticks":{"cutoff_price":"3.00","below_tick":"0.01","above_tick":"0.05"},"id":"74ecfc8e-3fee-4e70-85b6-d9fe755c96cc"},{"can_open_position":false,"symbol":"1SPY","trade_value_multiplier":"100.0000","underlying_instruments":[{"instrument":"https:\/\/api.robinhood.com\/instruments\/8f92e76f-1e0e-4478-8580-16a6ffcfaef5\/","id":"5b5a0dde-1f02-43a5-ac9a-5eb104dd380d","quantity":100}],"expiration_dates":[],"cash_component":null,"min_ticks":{"cutoff_price":"3.00","below_tick":"0.01","above_tick":"0.05"},"id":"de653940-25c0-4e35-986a-989737498881"}],"next":null}`,
	apiURL + optionsURI + "?chain_id=c277b118-58d9-4060-8dc5-a3b5898955cb&expiration_dates=2018-06-29&state=active&tradability=tradable": `{"previous":null,"results":[{"issue_date":"2005-01-06","tradability":"tradable","strike_price":"296.0000","url":"https:\/\/api.robinhood.com\/options\/instruments\/8ada9799-6c34-4647-b3ee-b6c157745740\/","expiration_date":"2018-06-29","created_at":"2018-06-02T10:16:57.966257Z","chain_id":"c277b118-58d9-4060-8dc5-a3b5898955cb","updated_at":"2018-06-02T10:16:57.966265Z","state":"active","type":"call","chain_symbol":"SPY","min_ticks":{"cutoff_price":"0.00","below_tick":"0.01","above_tick":"0.01"},"id":"8ada9799-6c34-4647-b3ee-b6c157745740"},{"issue_date":"2005-01-06","tradability":"tradable","strike_price":"298.0000","url":"https:\/\/api.robinhood.com\/options\/instruments\/637d839a-f3b3-45f9-91f4-b359c3ac80cb\/","expiration_date":"2018-06-29","created_at":"2018-06-02T10:16:57.964090Z","chain_id":"c277b118-58d9-4060-8dc5-a3b5898955cb","updated_at":"2018-06-02T10:16:57.964097Z","state":"active","type":"put","chain_symbol":"SPY","min_ticks":{"cutoff_price":"0.00","below_tick":"0.01","above_tick":"0.01"},"id":"637d839a-f3b3-45f9-91f4-b359c3ac80cb"}],"next":null}`,
	apiURL + oAuthUpgradeURI:                                            `{"token_type":"Bearer","access_token":"btok","expires_in":300,"refresh_token":"reftok","scope":"web_limited"}`,
	apiURL + marketOptionsURI + "8ada9799-6c34-4647-b3ee-b6c157745740/": `{"adjusted_mark_price":"28.3500","ask_price":"28.4700","ask_size":30,"bid_price":"28.2200","bid_size":35,"break_even_price":"269.6500","high_price":null,"instrument":"https://api.robinhood.com/options/instruments/637d839a-f3b3-45f9-91f4-b359c3ac80cb/","last_trade_price":null,"last_trade_size":null,"low_price":null,"mark_price":"28.3450","open_interest":0,"previous_close_date":"2018-06-22","previous_close_price":"23.2500","volume":0,"chance_of_profit_long":"0.5072","chance_of_profit_short":"0.4927","delta":"-0.9864","gamma":"0.0028","implied_volatility":"0.4268","rho":"-0.0322","theta":"-0.0400","vega":"0.0097"}`,
}

func TestOptions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for url, reply := range options {
		postOrGet := "GET"
		if strings.Contains(url, oAuthUpgradeURI) {
			postOrGet = "POST"
		}
		httpmock.RegisterResponder(postOrGet, url, httpmock.NewStringResponder(200, reply))
	}

	c := Client{
		AccountID: "account",
		Token:     "token",
	}
	exp, err := time.Parse("2006-01-02", "2018-06-29")
	if err != nil {
		t.Fatal(err)
	}
	chains, err := c.Chains("SPY", exp)
	if err != nil {
		t.Fatal(err)
	}
	if len(chains) < 1 {
		t.Fatalf("len(chains) = %d, want > 0", len(chains))
	}
	got, err := c.Option(chains[0])
	if err != nil {
		t.Fatal(err)
	}
	want := Option{"SPY", 296, exp, "call", 28.22, 28.47, 0, 28.35, 0, 0.4268, Instrument("https://api.robinhood.com/options/instruments/637d839a-f3b3-45f9-91f4-b359c3ac80cb/")}
	if got != want {
		t.Fatalf("want = %v, got = %v", want, got)
	}
}
