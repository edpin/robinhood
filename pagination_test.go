package robinhood

import (
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

var pages = map[string]string{
	"https://api.robinhood.com/options/":              `{"previous":null,"results":[{"result": "one"},{"result": "two"}], "next":"https://api.robinhood.com/options/?cursor=next1"}`,
	"https://api.robinhood.com/options/?cursor=next1": `{"previous":null,"results":[ {"result": "three"} ], "next":"https://api.robinhood.com/options/?cursor=next2"}`,
	"https://api.robinhood.com/options/?cursor=next2": `{"previous":null,"results":[], "next":"https://api.robinhood.com/options/?cursor=next3"}`,
	"https://api.robinhood.com/options/?cursor=next3": `{"previous":null,"results":[{"result": "four"}], "next":null}`,
}

func TestPagination(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for url, reply := range pages {
		httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, reply))
	}

	c := Client{
		AccountID: "account",
		Token:     "token",
	}
	resp, err := c.paginatedGet("options/")
	if err != nil {
		t.Fatal(err)
	}
	// Check results.
	want := `[{"result": "one"},{"result": "two"},{"result": "three"},{"result": "four"}]`
	if string(resp) != want {
		t.Fatalf("got = %q, want = %q", resp, want)
	}
}
