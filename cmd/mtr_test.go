package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInitMTRCmd(t *testing.T) {
	parent := &cobra.Command{}
	initMTRCmd(parent)

	flags := mtrCmd.Flags()

	if got := flags.Lookup("limit"); got == nil {
		t.Fatal("expected limit flag; got nil")
	}
}

func TestRunMTR(t *testing.T) {
	// We're only interested in the first HTTP call, e.g., the one to get the test ID
	// to validate our parameters got passed properly.
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	runMTR(c, "example.com", "From here", 123)
	if got, exp := tr.req.URL.Path, "/run/mtr"; got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
	got := reqBody(tr.req)
	exp := `{"target":"example.com","location":"From here","limit":123}`
	if got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
}
