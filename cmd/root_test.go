// Copyright 2017 The PerfOps-CLI Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	"github.com/spf13/cobra"

	"github.com/ProspectOne/perfops-cli/perfops"
)

func TestEnvPerfOpsAPIKey(t *testing.T) {
	os.Unsetenv("PERFOPS_API_KEY")

	rootCmd.ResetFlags()
	initRootCmd()
	if err := rootCmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != "" {
		t.Fatalf("expected no key; got %v", got)
	}

	const apiKey2 = "Moo"
	rootCmd.ResetFlags()
	initRootCmd()
	if err := rootCmd.ParseFlags([]string{"-K", apiKey2}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != apiKey2 {
		t.Fatalf("expected %v; got %v", apiKey2, got)
	}

	const envAPIKey = "Meep"
	os.Setenv("PERFOPS_API_KEY", envAPIKey)
	defer os.Unsetenv("PERFOPS_API_KEY")
	rootCmd.ResetFlags()
	initRootCmd()
	if err := rootCmd.ParseFlags([]string{}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got := apiKey; got != envAPIKey {
		t.Fatalf("expected %v; got %v", envAPIKey, got)
	}

	rootCmd.ResetFlags()
	initRootCmd()
	if err := rootCmd.ParseFlags([]string{"-K", apiKey2}); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if err := rootCmd.Execute(); err != nil {
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

func TestShowErrorOnly(t *testing.T) {
	expErr := errors.New("error")
	cmd := &cobra.Command{
		Use: "test-only",
		RunE: func(cmd *cobra.Command, args []string) error {
			return expErr
		},
	}
	var b bytes.Buffer
	rootCmd.ResetFlags()
	rootCmd.SetOutput(&b)
	rootCmd.AddCommand(cmd)
	initRootCmd()
	rootCmd.SetArgs([]string{"test-only"})
	if err := rootCmd.Execute(); err != expErr {
		t.Fatalf("exepected %v; got %v", expErr, err)
	}
	if got, exp := b.String(), "Error: error\n"; got != exp {
		t.Fatalf("expected %#v; got %q", exp, got)
	}
}

type unauthedError struct{}

func (e *unauthedError) Error() string        { return "unauthorized" }
func (e *unauthedError) IsUnauthorized() bool { return true }

type argError struct {
	name string
}

func (e *argError) Error() string   { return fmt.Sprintf("invalid argument: %s", e.name) }
func (e *argError) ArgName() string { return e.name }

func TestChkRunError(t *testing.T) {
	testCases := map[string]struct {
		err error
		exp error
	}{
		"Param error":   {&argError{"flag"}, errors.New("pflag: help requested")},
		"401":           {&unauthedError{}, errors.New("The API token was declined. Please correct it or do not send a token to use the free plan.")},
		"Generic error": {errors.New("meep"), errors.New("meep")},
		"No error":      {nil, nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if got, exp := chkRunError(tc.err), tc.exp; !cmpError(got, exp) {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}

func TestInvalidArgHelp(t *testing.T) {
	var fo string
	var fr string
	var fl string
	var b bytes.Buffer
	cmd := &cobra.Command{}
	cmd.Flags().StringVarP(&fo, "flag", "F", "", "A flag")
	cmd.Flags().StringVarP(&fr, "req", "R", "", "A second flag")
	cmd.Flags().StringVarP(&fl, "limit", "L", "", "Limit")
	cmd.MarkFlagRequired("req")
	cmd.SetHelpFunc(invalidArgHelp("flag"))
	cmd.SetOutput(&b)
	if err := cmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got, exp := b.String(), "Missing or invalid arguments:\n  -F, --flag string\n  -R, --req string\n"; got != exp {
		t.Fatalf("expected %q; got %q", exp, got)
	}

	b.Reset()
	cmd.SetHelpFunc(invalidArgHelp("limit"))
	if err := cmd.Execute(); err != nil {
		t.Fatalf("exepected nil; got %v", err)
	}
	if got, exp := b.String(), "For free users the maximum allowed number nodes for a single test is 20. Please change your limit.\n"; got != exp {
		t.Fatalf("expected %q; got %q", exp, got)
	}
}

func cmpError(a, b error) bool {
	return a == b || (a != nil && b != nil && a.Error() == b.Error())
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
