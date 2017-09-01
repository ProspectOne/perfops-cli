package perfops

import (
	"context"
	"testing"
)

func TestPing(t *testing.T) {
	testCases := map[string]struct {
		pingID string
		tr     *respondingTransport
		err    bool
	}{
		"HTTP error": {"", &respondingTransport{resp: dummyResp(400, "POST", `{"Error": "an error"}`)}, true},
		"Failed":     {"", &respondingTransport{resp: dummyResp(201, "POST", `{"Error": "an error"}`)}, true},
		"Created":    {"0123456789abcdefghij", &respondingTransport{resp: dummyResp(201, "POST", `{"id": "0123456789abcdefghij"}`)}, false},
	}
	ctx := context.Background()
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, err := newTestClient(tc.tr)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := c.Run.Ping(ctx, &Ping{})
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
			got, err := c.Run.PingOutput(ctx, PingID("1234"))
			if (err == nil && tc.err) || (err != nil && !tc.err) {
				t.Fatalf("expected error %v; got %v", tc.err, err)
			}
			if got.IsFinished() != tc.finished {
				t.Fatalf("expected %v; got %v", tc.finished, got.IsFinished())
			}
		})
	}
}
