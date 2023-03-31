// Copyright 2017 Prospect One https://prospectone.io/. All rights reserved.
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
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitDNSPerfCmd(t *testing.T) {
	testCases := map[string]struct {
		args   []string
		gotexp func() (interface{}, interface{})
	}{
		// Common flags
		"from":   {[]string{"--from", "Europe"}, func() (interface{}, interface{}) { return from, "Europe" }},
		"nodeid": {[]string{"--nodeid", "1,2,3"}, func() (interface{}, interface{}) { return nodeIDs, []int{1, 2, 3} }},
		"json":   {[]string{"--json"}, func() (interface{}, interface{}) { return outputJSON, true }},

		"dns-server": {[]string{"--dns-server", "123.234.0.1"}, func() (interface{}, interface{}) { return dnsPerfDNSServer, "123.234.0.1" }},
		"limit":      {[]string{"--limit", "23"}, func() (interface{}, interface{}) { return dnsPerfLimit, 23 }},
		"ipv6":       {[]string{"--ipv6"}, func() (interface{}, interface{}) { return dnsPerfIpv6, true }},
	}
	parent := &cobra.Command{}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			dnsPerfCmd.ResetFlags()
			initDNSPerfCmd(parent)
			if err := dnsPerfCmd.ParseFlags(tc.args); err != nil {
				t.Fatalf("exepected nil; got %v", err)
			}
			flags := dnsPerfCmd.Flags()
			f := flags.Lookup(name)
			if f == nil {
				t.Fatal("expected flag; got nil")
			}

			got, exp := tc.gotexp()
			if reflect.DeepEqual(got, exp) == false {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}

func TestRunDNSPerf(t *testing.T) {
	testCases := map[string]struct {
		from    string
		nodeIDs []int
		ipv6    bool
		exp     string
	}{
		"Location": {"From here", []int{}, false, `{"target":"example.com","dnsServer":"127.0.0.1","location":"From here","limit":12,"ipversion":4}`},
		"NodeID":   {"", []int{123}, false, `{"target":"example.com","dnsServer":"127.0.0.1","nodes":"123","limit":12,"ipversion":4}`},
		"IPv6":     {"", []int{}, true, `{"target":"example.com","dnsServer":"127.0.0.1","limit":12,"ipversion":6}`},
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
			runDNSPerf(c, "example.com", "127.0.0.1", tc.from, tc.nodeIDs, 12, tc.ipv6)
			if got, exp := tr.req.URL.Path, "/run/dns-perf"; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
			if got, exp := reqBody(tr.req), tc.exp; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}
