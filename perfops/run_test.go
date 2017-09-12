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

package perfops

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLatency(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.Latency(ctx, &RunRequest{Target: "example.com"})
	if got, exp := tr.req.Method, "POST"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/latency"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestLatencyOutput(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.LatencyOutput(ctx, TestID("1234"))
	if got, exp := tr.req.Method, "GET"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/latency/1234"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}

	const body = `{"id": "2e9fd0e3a444adddb9b8168e6e0f856c","items": [{"id": "186b4c4c77985f75e7cefc48289e79ff","result": {"ip": "74.125.200.113","output": "35.223","node": {"id": 58,"latitude": 22.280042521009999,"longitude": 114.1915512085,"country": {"id": 195,"name": "Hong Kong","continent": {"id": 2,"name": "Asia","iso": "AS"},"iso": "HK","iso_numeric": "344"},"city": "Hong Kong","sub_region": "Eastern Asia"}}}, {"id": "34ae843cf4d341b252a14e3d6b39281f","result": {"output": -1,"message": "100% packet loss","node": {"id": 208,"latitude": 30.255360295637001,"longitude": 120.15712738037,"country": {"id": 62,"name": "China","continent": {"id": 2,"name": "Asia","iso": "AS"},"iso": "CN","iso_numeric": "156"},"city": "Hangzhou","sub_region": "Eastern Asia"}}}],"requested": "google.com","finished": "true"}`
	tr2 := &respondingTransport{resp: dummyResp(201, "GET", body)}
	c, err = newTestClient(tr2)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	_, err = c.Run.LatencyOutput(ctx, TestID("e5c97cee7806ca44577cc92308d3b8e8"))
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}

func TestMTR(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.MTR(ctx, &RunRequest{Target: "example.com"})
	if got, exp := tr.req.Method, "POST"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/mtr"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestMTROutput(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.MTROutput(ctx, TestID("1234"))
	if got, exp := tr.req.Method, "GET"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/mtr/1234"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestPing(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.Ping(ctx, &RunRequest{Target: "example.com"})
	if got, exp := tr.req.Method, "POST"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/ping"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestPingOutput(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.PingOutput(ctx, TestID("1234"))
	if got, exp := tr.req.Method, "GET"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/ping/1234"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestTraceroute(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.Traceroute(ctx, &RunRequest{Target: "example.com"})
	if got, exp := tr.req.Method, "POST"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/traceroute"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestTracerouteOutput(t *testing.T) {
	ctx := context.Background()
	tr := &recordingTransport{}
	c, err := newTestClient(tr)
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	c.Run.TracerouteOutput(ctx, TestID("1234"))
	if got, exp := tr.req.Method, "GET"; got != exp {
		t.Fatalf("expected HTTP method %v; got %v", exp, got)
	}
	if got, exp := tr.req.URL.Path, "/run/traceroute/1234"; got != exp {
		t.Fatalf("expected path %v; got %v", exp, got)
	}
}

func TestDNSResolve(t *testing.T) {
	reqTestCases := map[string]struct {
		dnsResolveReq DNSResolveRequest
		tr            *recordingTransport
		expReqBody    string
	}{
		"Required only": {DNSResolveRequest{Target: "example.com", Param: "A", DNSServer: "127.0.0.1"}, &recordingTransport{}, `{"target":"example.com","param":"A","dnsServer":"127.0.0.1"}`},
		"With location": {DNSResolveRequest{Target: "example.com", Param: "A", DNSServer: "127.0.0.1", Location: "Asia"}, &recordingTransport{}, `{"target":"example.com","param":"A","dnsServer":"127.0.0.1","location":"Asia"}`},
	}
	ctx := context.Background()
	for name, tc := range reqTestCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			c.Run.DNSResolve(ctx, &tc.dnsResolveReq)
			if got := reqBody(tc.tr.req); tc.expReqBody != "" && tc.expReqBody != got {
				t.Fatalf("expected %v; got %v", tc.expReqBody, got)
			}
		})
	}

	testCases := map[string]struct {
		target    string
		param     string
		dnsServer string
		limit     int
		testID    string
		tr        *respondingTransport
		err       error
	}{
		"Invalid target":     {"meep", "A", "127.0.0.1", 1, "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, &argError{"target"}},
		"Invalid param":      {"example.com", "", "127.0.0.1", 1, "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, &argError{"param"}},
		"Missing DNS server": {"example.com", "A", "", 1, "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, &argError{"dns server"}},
		"Invalid DNS server": {"example.com", "A", "127.0", 1, "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, &argError{"dns server"}},
		"Invalid limit":      {"example.com", "A", "127.0.0.1", freeMaxNodeCap + 1, "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, &argError{"limit"}},
		"HTTP error":         {"example.com", "A", "127.0.0.1", 1, "", &respondingTransport{resp: dummyResp(400, "POST", `{"Error": "an error"}`)}, errors.New(`400: {"Error": "an error"}`)},
		"Failed":             {"example.com", "A", "127.0.0.1", 1, "", &respondingTransport{resp: dummyResp(201, "POST", `{"Error": "an error"}`)}, errors.New("an error")},
		"Created":            {"example.com", "A", "127.0.0.1", 1, "0123456789abcdefghij", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "0123456789abcdefghij"}`)}, nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.DNSResolve(ctx, &DNSResolveRequest{Target: tc.target, Param: tc.param, DNSServer: tc.dnsServer, Limit: tc.limit})
			if !cmpError(err, tc.err) {
				t.Fatalf("expected %v; got %v", tc.err, err)
			}
			if string(got) != tc.testID {
				t.Fatalf("expected %v; got %v", tc.testID, got)
			}
		})
	}
}

func TestDNSResolveOutput(t *testing.T) {
	testCases := map[string]struct {
		tr       *respondingTransport
		err      error
		finished bool
	}{
		"Incomplete": {&respondingTransport{resp: dummyResp(200, "GET", `{"id":"d1f2408ff","items":[{"id":"734df82","result":{"id":123,"message":"NO DATA"}}]}`)}, nil, false},
		"Complete":   {&respondingTransport{resp: dummyResp(200, "GET", `{"id": "66b78cfc643ea238e0fd8ab44f512657","items": [{"id": "ae3e8bcd0fbe77d6322b89371d87d96d","result": {"dnsServer": "8.8.8.8","output": ["204.79.197.200","13.107.21.200"],"node": {"id": 5,"latitude": 50.110781326572834,"longitude": 8.68984222412098,"country": {"id": 116,"name": "Germany","continent": {"id": 3,"name": "Europe","iso": "EU"},"iso": "DE","iso_numeric": "276"},"city": "Frankfurt","sub_region": "Western Europe"}}}],"requested": "bing.com","finished": "true"}`)}, nil, true},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.DNSResolveOutput(ctx, TestID("1234"))
			if !cmpError(err, tc.err) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if got.IsFinished() != tc.finished {
				t.Fatalf("expected %v; got %v", tc.finished, got.IsFinished())
			}
		})
	}
}

func TestIsValidTarget(t *testing.T) {
	testCases := map[string]struct {
		t     string
		valid bool
	}{
		"Empty target":     {"", false},
		"Invalid hostname": {"meep", false},
		"Valid hostname":   {"meep.com", true},
		"Valid IPv4":       {"123.123.123.123", true},
		"Invalid IPv6":     {"123:123", false},
		"Valid IPv6":       {"2001:db8::68", true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got := isValidTarget(tc.t)
			if got != tc.valid {
				t.Fatalf("expected %v; got %v", tc.valid, got)
			}
		})
	}
}

func TestIsValidLimit(t *testing.T) {
	testCases := map[string]struct {
		key   string
		limit int
		exp   bool
	}{
		"No key, limit 1":   {"", 1, true},
		"No key, limit 10":  {"", 10, true},
		"No key, limit 20":  {"", 20, true},
		"No key, limit 21":  {"", 21, false},
		"No key, limit 221": {"", 221, false},
		"Key, limit 1":      {"key", 1, true},
		"Key, limit 10":     {"key", 10, true},
		"Key, limit 20":     {"key", 20, true},
		"Key, limit 21":     {"key", 21, true},
		"Key, limit 221":    {"key", 221, true},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if got, exp := isValidLimit(tc.key, tc.limit), tc.exp; got != exp {
				t.Fatalf("expected %v; got %v", exp, got)
			}
		})
	}
}

func TestDoPostRunRequest(t *testing.T) {
	errDummyTr := errors.New("Post https://api.perfops.net/run/test: dummy impl")
	reqTestCases := map[string]struct {
		runReq     RunRequest
		tr         *recordingTransport
		expErr     error
		expReqBody string
	}{
		"Invalid target":          {RunRequest{Target: "example"}, &recordingTransport{}, &argError{"target"}, ``},
		"Target only":             {RunRequest{Target: "example.com"}, &recordingTransport{}, errDummyTr, `{"target":"example.com"}`},
		"With invalid limit":      {RunRequest{Target: "example.com", Limit: freeMaxNodeCap + 2}, &recordingTransport{}, &argError{"limit"}, ``},
		"With limit":              {RunRequest{Target: "example.com", Limit: 2}, &recordingTransport{}, errDummyTr, `{"target":"example.com","limit":2}`},
		"With location":           {RunRequest{Target: "example.com", Location: "Asia"}, &recordingTransport{}, errDummyTr, `{"target":"example.com","location":"Asia"}`},
		"With limit and location": {RunRequest{Target: "example.com", Limit: 2, Location: "Asia"}, &recordingTransport{}, errDummyTr, `{"target":"example.com","location":"Asia","limit":2}`},
	}
	ctx := context.Background()
	for name, tc := range reqTestCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			_, err = c.Run.doPostRunRequest(ctx, "/run/test", &tc.runReq)
			if !cmpError(err, tc.expErr) {
				t.Fatalf("expected %v; got %v", tc.expErr, err)
			}
			if got := reqBody(tc.tr.req); tc.expReqBody != "" && tc.expReqBody != got {
				t.Fatalf("expected %v; got %v", tc.expReqBody, got)
			}
		})
	}

	respTestCases := map[string]struct {
		target string
		testID string
		tr     *respondingTransport
		err    error
	}{
		"Invalid target": {"meep", "", &respondingTransport{}, &argError{"target"}},
		"HTTP error":     {"example.com", "", &respondingTransport{resp: dummyResp(400, "POST", `{"Error": "an error"}`)}, errors.New(`400: {"Error": "an error"}`)},
		"Unauthorized":   {"example.com", "", &respondingTransport{resp: dummyResp(401, "POST", `Unauthorized`)}, errors.New(`401: Unauthorized`)},
		"Failed":         {"example.com", "", &respondingTransport{resp: dummyResp(201, "POST", `{"Error": "an error"}`)}, errors.New("an error")},
		"Created":        {"example.com", "0123456789abcdefghij", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "0123456789abcdefghij"}`)}, nil},
	}
	for name, tc := range respTestCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.doPostRunRequest(ctx, "/run/test", &RunRequest{Target: tc.target})
			if !cmpError(err, tc.err) {
				t.Fatalf("expected %v; got %v", tc.err, err)
			}
			if string(got) != tc.testID {
				t.Fatalf("expected %v; got %v", tc.testID, got)
			}
		})
	}
}

func TestDoGetRunOutput(t *testing.T) {
	testCases := map[string]struct {
		tr       *respondingTransport
		err      error
		finished bool
	}{
		"Incomplete": {&respondingTransport{resp: dummyResp(200, "GET", `{"id":"d1f2408ff","items":[{"id":"734df82","result":{"id":123,"message":"NO DATA"}}]}`)}, nil, false},
		"Complete":   {&respondingTransport{resp: dummyResp(200, "GET", `{"id": "65d2bb722be16277e3fa8e8c86d3afb7","items": [{"id": "0981fcaf99f2c1b4a46a22cedb417347","result": {"output": "PING bing.com (204.79.197.200): 56 data bytes\n64 bytes from 204.79.197.200: icmp_seq=0 ttl=119 time=40.348 ms\n64 bytes from 204.79.197.200: icmp_seq=1 ttl=119 time=40.198 ms\n64 bytes from 204.79.197.200: icmp_seq=2 ttl=119 time=40.241 ms\n--- bing.com ping statistics ---\n3 packets transmitted, 3 packets received, 0% packet loss\nround-trip min/avg/max/stddev = 40.198/40.262/40.348/0.063 ms\n","node": {"id": 5,"latitude": 50.110781326572834,"longitude": 8.68984222412098,"country": {"id": 116,"name": "Germany","continent": {"id": 3,"name": "Europe","iso": "EU"},"iso": "DE","iso_numeric": "276"},"city": "Frankfurt","sub_region": "Western Europe"}}}],"requested": "bing.com","finished": "true"}`)}, nil, true},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.PingOutput(ctx, TestID("1234"))
			if !cmpError(err, tc.err) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if got.IsFinished() != tc.finished {
				t.Fatalf("expected %v; got %v", tc.finished, got.IsFinished())
			}
		})
	}
}

func cmpError(a, b error) bool {
	return a == b || (a != nil && b != nil && a.Error() == b.Error())
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
