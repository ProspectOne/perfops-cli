package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInitPingCmd(t *testing.T) {
	parent := &cobra.Command{}
	initPingCmd(parent)

	flags := pingCmd.Flags()

	if got := flags.Lookup("limit"); got == nil {
		t.Fatal("expected limit flag; got nil")
	}
}

func TestRunPing(t *testing.T) {
	// We're only interested in the first HTTP call, e.g., the one to get the test ID
	// to validate our parameters got passed properly.
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	runPing(c, "example.com", "From here", 123)
	if got, exp := tr.req.URL.Path, "/run/ping"; got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
	got := reqBody(tr.req)
	exp := `{"target":"example.com","location":"From here","limit":123}`
	if got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
}
