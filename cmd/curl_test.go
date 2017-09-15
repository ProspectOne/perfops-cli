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
	"testing"

	"github.com/spf13/cobra"
)

func TestInitCUrlCmd(t *testing.T) {
	parent := &cobra.Command{}
	curlCmd.ResetFlags()
	initCurlCmd(parent)

	flags := curlCmd.Flags()

	if got := flags.Lookup("head"); got == nil {
		t.Fatal("expected head flag; got nil")
	}
	if got := flags.Lookup("insecure"); got == nil {
		t.Fatal("expected insecure flag; got nil")
	}
	if got := flags.Lookup("http2"); got == nil {
		t.Fatal("expected http2 flag; got nil")
	}
}

func TestRunCurlResolve(t *testing.T) {
	testCases := map[string]struct {
		head     bool
		insecure bool
		http2    bool
		from     string
		nodeIDs  []int
		exp      string
	}{
		"Head":     {false, false, false, "From here", []int{}, `{"target":"example.com","head":false,"location":"From here","limit":12}`},
		"Insecure": {true, true, false, "From here", []int{}, `{"target":"example.com","head":true,"insecure":true,"location":"From here","limit":12}`},
		"HTTP2":    {true, false, true, "From here", []int{}, `{"target":"example.com","head":true,"http2":true,"location":"From here","limit":12}`},
		"Location": {true, false, false, "From here", []int{}, `{"target":"example.com","head":true,"location":"From here","limit":12}`},
		"NodeID":   {true, false, false, "", []int{123}, `{"target":"example.com","head":true,"nodes":"123","limit":12}`},
	}
	// We're only interested in the first HTTP call, e.g., the one to get the test ID
	// to validate our parameters got passed properly.
	tr := &recordingTransport{}
	c, err := newTestPerfopsClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			runCurl(c, "example.com", tc.head, tc.insecure, tc.http2, tc.from, tc.nodeIDs, 12)
			if got, exp := tr.req.URL.Path, "/run/curl"; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
			if got, exp := reqBody(tr.req), tc.exp; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}
