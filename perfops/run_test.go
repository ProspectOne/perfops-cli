package perfops

import (
	"context"
	"testing"
)

func TestMTR(t *testing.T) {
	testCases := map[string]struct {
		target string
		mtrID  string
		tr     *respondingTransport
		err    bool
	}{
		"Invalid target": {"meep", "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, true},
		"HTTP error":     {"example.com", "", &respondingTransport{resp: dummyResp(400, "POST", `{"Error": "an error"}`)}, true},
		"Failed":         {"example.com", "", &respondingTransport{resp: dummyResp(201, "POST", `{"Error": "an error"}`)}, true},
		"Created":        {"example.com", "0123456789abcdefghij", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "0123456789abcdefghij"}`)}, false},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.MTR(ctx, &RunRequest{Target: tc.target})
			if (err == nil && tc.err) || (err != nil && !tc.err) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if string(got) != tc.mtrID {
				t.Fatalf("expected %v; got %v", tc.mtrID, got)
			}
		})
	}
}

func TestMTROutput(t *testing.T) {
	testCases := map[string]struct {
		tr       *respondingTransport
		err      bool
		finished bool
	}{
		"Incomplete": {&respondingTransport{resp: dummyResp(200, "GET", `{"id":"d1f2408ff","items":[{"id":"734df82","result":{"id":123,"message":"NO DATA"}}]}`)}, false, false},
		"Complete":   {&respondingTransport{resp: dummyResp(200, "GET", `{"id": "9072a72f762b876525ca4c9153af9983","items": [{"id": "edca088e43bde5453b961f6210723157","result": {"output": "Start: Thu Jul 27 15:59:05 2017                Loss%   Snt   Last   Avg  Best  Wrst StDev\n  1.|-- 172.18.0.1                 0.0%     2    0.0   0.1   0.0   0.1   0.0\n  2.|-- 10.0.2.2                   0.0%     2    0.2   0.2   0.2   0.2   0.0\n  3.|-- 192.168.0.1                0.0%     2    1.3   1.5   1.3   1.6   0.0\n  4.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n  5.|-- 80.81.194.168              0.0%     2   25.0  23.4  21.8  25.0   2.0\n  6.|-- 80.81.194.52               0.0%     2   41.4  41.4  41.4  41.4   0.0\n  7.|-- 104.44.80.143              0.0%     2   40.9  40.5  40.0  40.9   0.0\n  8.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n  9.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n 10.|-- ???                       100.0     2    0.0   0.0   0.0   0.0   0.0\n 11.|-- 13.107.21.200              0.0%     2   39.8  40.1  39.8  40.3   0.0\n","node": {"id": 5,"latitude": 50.110781326572834,"longitude": 8.68984222412098,"country": {"id": 116,"name": "Germany","continent": {"id": 3,"name": "Europe","iso": "EU"},"iso": "DE","iso_numeric": "276"},"city": "Frankfurt","sub_region": "Western Europe"}}}],"requested": "bing.com","finished": "true"}`)}, false, true},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.MTROutput(ctx, TestID("1234"))
			if (err == nil && tc.err) || (err != nil && !tc.err) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if got.IsFinished() != tc.finished {
				t.Fatalf("expected %v; got %v", tc.finished, got.IsFinished())
			}
		})
	}
}

func TestPing(t *testing.T) {
	testCases := map[string]struct {
		target string
		pingID string
		tr     *respondingTransport
		err    bool
	}{
		"Invalid target": {"meep", "", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "135"}`)}, true},
		"HTTP error":     {"example.com", "", &respondingTransport{resp: dummyResp(400, "POST", `{"Error": "an error"}`)}, true},
		"Failed":         {"example.com", "", &respondingTransport{resp: dummyResp(201, "POST", `{"Error": "an error"}`)}, true},
		"Created":        {"example.com", "0123456789abcdefghij", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "0123456789abcdefghij"}`)}, false},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.Ping(ctx, &RunRequest{Target: tc.target})
			if (err == nil && tc.err) || (err != nil && !tc.err) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if string(got) != tc.pingID {
				t.Fatalf("expected %v; got %v", tc.pingID, got)
			}
		})
	}
}

func TestPingOutput(t *testing.T) {
	testCases := map[string]struct {
		tr       *respondingTransport
		err      bool
		finished bool
	}{
		"Incomplete": {&respondingTransport{resp: dummyResp(200, "GET", `{"id":"d1f2408ff","items":[{"id":"734df82","result":{"id":123,"message":"NO DATA"}}]}`)}, false, false},
		"Complete":   {&respondingTransport{resp: dummyResp(200, "GET", `{"id": "65d2bb722be16277e3fa8e8c86d3afb7","items": [{"id": "0981fcaf99f2c1b4a46a22cedb417347","result": {"output": "PING bing.com (204.79.197.200): 56 data bytes\n64 bytes from 204.79.197.200: icmp_seq=0 ttl=119 time=40.348 ms\n64 bytes from 204.79.197.200: icmp_seq=1 ttl=119 time=40.198 ms\n64 bytes from 204.79.197.200: icmp_seq=2 ttl=119 time=40.241 ms\n--- bing.com ping statistics ---\n3 packets transmitted, 3 packets received, 0% packet loss\nround-trip min/avg/max/stddev = 40.198/40.262/40.348/0.063 ms\n","node": {"id": 5,"latitude": 50.110781326572834,"longitude": 8.68984222412098,"country": {"id": 116,"name": "Germany","continent": {"id": 3,"name": "Europe","iso": "EU"},"iso": "DE","iso_numeric": "276"},"city": "Frankfurt","sub_region": "Western Europe"}}}],"requested": "bing.com","finished": "true"}`)}, false, true},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.PingOutput(ctx, TestID("1234"))
			if (err == nil && tc.err) || (err != nil && !tc.err) {
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
