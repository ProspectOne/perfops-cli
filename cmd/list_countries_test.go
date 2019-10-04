package cmd

import (
	"testing"
)

func TestGetCountriesList(t *testing.T) {
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	runCountriesCmd(c)
	if got, exp := tr.req.URL.Path, "/analytics/dns/countries"; got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
}
