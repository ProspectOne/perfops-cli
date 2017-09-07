package perfops

import (
	"context"
	"errors"
	"fmt"
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

func TestDoPostRunRequest(t *testing.T) {
	testCases := map[string]struct {
		target string
		testID string
		tr     *respondingTransport
		err    error
	}{
		"Invalid target": {"meep", "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, errors.New("target invalid")},
		"HTTP error":     {"example.com", "", &respondingTransport{resp: dummyResp(400, "POST", `{"Error": "an error"}`)}, fmt.Errorf("HTTP Error: %v", http.StatusBadRequest)},
		"Failed":         {"example.com", "", &respondingTransport{resp: dummyResp(201, "POST", `{"Error": "an error"}`)}, errors.New("an error")},
		"Created":        {"example.com", "0123456789abcdefghij", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "0123456789abcdefghij"}`)}, nil},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.doPostRunRequest(ctx, "/run/test", &RunRequest{Target: tc.target})
			if (err == nil && tc.err != nil) || (err != nil && tc.err == nil) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
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
			if (err == nil && tc.err != nil) || (err != nil && tc.err == nil) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if got.IsFinished() != tc.finished {
				t.Fatalf("expected %v; got %v", tc.finished, got.IsFinished())
			}
		})
	}
}
