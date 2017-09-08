package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestInitDNSResolveCmd(t *testing.T) {
	parent := &cobra.Command{}
	initDNSResolveCmd(parent)

	flags := dnsResolveCmd.Flags()

	if got := flags.Lookup("param"); got == nil {
		t.Fatal("expected param flag; got nil")
	}
	if got := flags.Lookup("dns-server"); got == nil {
		t.Fatal("expected dns-server flag; got nil")
	}
}

func TestRunDNSResolve(t *testing.T) {
	// We're only interested in the first HTTP call, e.g., the one to get the test ID
	// to validate our parameters got passed properly.
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	runDNSResolve(c, "example.com", "TXT", "127.0.0.1", "From here")
	if got, exp := tr.req.URL.Path, "/run/dns-resolve"; got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
	got := reqBody(tr.req)
	exp := `{"target":"example.com","param":"TXT","dnsServer":"127.0.0.1","location":"From here"}`
	if got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
}
