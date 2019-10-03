package cmd

import (
	"testing"
)

func TestGetCitiesList(t *testing.T) {
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}

	runCitiesCmd(c)
	if got, exp := tr.req.URL.Path, "/analytics/dns/city"; got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
}
