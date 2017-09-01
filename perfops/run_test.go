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
		complete bool
	}{
		"Incomplete": {&respondingTransport{resp: dummyResp(200, "GET", `{"id":"d1f2408ff","items":[{"id":"734df82","result":{"id":123,"message":"NO DATA"}}]}`)}, false, false},
		"Complete":   {&respondingTransport{resp: dummyResp(200, "GET", `{"id": "0eb43618d","items": [{"id": "a01b51d5f","result": {"output": "PING example.com (13.107.21.200) 56(84) bytes of data.\n64 bytes from 13.107.21.200: icmp_seq=1 ttl=119 time=3.84 ms\n64 bytes from 13.107.21.200: icmp_seq=2 ttl=119 time=3.78 ms\n64 bytes from 13.107.21.200: icmp_seq=3 ttl=119 time=4.10 ms\n\n--- bing.com ping statistics ---\n3 packets transmitted, 3 received, 0% packet loss, time 601ms\nrtt min/avg/max/mdev = 3.781/3.909/4.108/0.159 ms\n","nodeId": "123"}},{"id": "76dd05a41d22c835763bb1fb7c0f00cb","result": {"output": "PING example.com (204.79.197.200) 56(84) bytes of data.\n64 bytes from 204.79.197.200: icmp_seq=1 ttl=122 time=3.88 ms\n64 bytes from 204.79.197.200: icmp_seq=2 ttl=122 time=4.22 ms\n64 bytes from 204.79.197.200: icmp_seq=3 ttl=122 time=4.08 ms\n\n--- bing.com ping statistics ---\n3 packets transmitted, 3 received, 0% packet loss, time 603ms\nrtt min/avg/max/mdev = 3.881/4.062/4.223/0.149 ms\n","nodeId": "124"}}],"requested": "example.com"}`)}, false, true},
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
			if got.IsComplete() != tc.complete {
				t.Fatalf("expected %v; got %v", tc.complete, got.IsComplete())
			}
		})
	}
}
