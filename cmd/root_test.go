package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"testing"

	"github.com/ProspectOne/perfops-cli/perfops"
)

func TestEnvPerfOpsAPIKey(t *testing.T) {
	os.Unsetenv("PERFOPS_API_KEY")

	rootCmd.ResetFlags()
	initRootCmd()

	if err := rootCmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != "" {
		t.Fatalf("expected no key; got %v", got)
	}

	const envAPIKey = "Meep"
	os.Setenv("PERFOPS_API_KEY", envAPIKey)
	defer os.Unsetenv("PERFOPS_API_KEY")

	rootCmd.ResetFlags()
	initRootCmd()

	if got := apiKey; got != envAPIKey {
		t.Fatalf("expected %v; got %v", envAPIKey, got)
	}

	const apiKey2 = "Moo"
	if err := rootCmd.ParseFlags([]string{"-K", apiKey2}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != apiKey2 {
		t.Fatalf("expected %v; got %v", apiKey2, got)
	}
}

func TestVersionFlag(t *testing.T) {
	versionOutput := fmt.Sprintf(versionTmpl, version, buildDate, commitHash,
		runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	var b bytes.Buffer
	rootCmd.ResetFlags()
	rootCmd.SetOutput(&b)
	initRootCmd()

	if err := rootCmd.ParseFlags([]string{"-v"}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got, exp := showVersion, true; got != exp {
		t.Fatalf("expected %v; got %v", exp, got)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got, exp := b.String(), versionOutput; got != exp {
		t.Fatalf("expected %q; got %q", exp, got)
	}
}

func TestUsage(t *testing.T) {
	var b bytes.Buffer
	rootCmd.ResetFlags()
	rootCmd.SetOutput(&b)
	initRootCmd()

	if err := rootCmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := b.String(); got == "" {
		t.Fatalf("expected not empty string; got %q", got)
	}
}

type roundTripper interface {
	RoundTrip(req *http.Request) (*http.Response, error)
}

func newTestPerfopsClient(tr roundTripper) (*perfops.Client, error) {
	c := &http.Client{Transport: tr}
	return perfops.NewClient(perfops.WithHTTPClient(c))
}

type recordingTransport struct {
	req *http.Request
}

func (t *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.req = req
	return nil, errors.New("dummy impl")
}

func reqBody(req *http.Request) string {
	if req == nil || req.Body == nil {
		return ""
	}

	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return ""
	}
	return string(bytes.TrimSpace(b))
}
